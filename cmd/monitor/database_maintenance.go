package main

import (
 "context"
 "database/sql"
 "fmt"
 "log"
 "os"
 "time"
)

const (
 walWarningBytes int64 = 256 << 20
 walCriticalBytes int64 = 1 << 30
)

// pruneOldSamples keeps the existing 30-day retention policy but runs daily,
// rather than issuing a DELETE every 30 seconds.
func pruneOldSamples(ctx context.Context, db *sql.DB, now time.Time) (int64,error) {
 result,err:=db.ExecContext(ctx,"DELETE FROM samples WHERE ts < ?",now.Add(-30*24*time.Hour).Unix())
 if err!=nil{return 0,err}
 return result.RowsAffected()
}

// A successful SQLite checkpoint has busy==0. Never silently delete a WAL.
func checkpointWAL(ctx context.Context, db *sql.DB)(busy,frames,checkpointed int,err error){
 err=db.QueryRowContext(ctx,"PRAGMA wal_checkpoint(TRUNCATE)").Scan(&busy,&frames,&checkpointed)
 return
}
func walSize(path string)(int64,error){
 st,err:=os.Stat(path+"-wal")
 if os.IsNotExist(err){return 0,nil}
 if err!=nil{return 0,err}
 return st.Size(),nil
}
func maybeCheckpointWAL(ctx context.Context,db *sql.DB,path string,threshold int64) error{
 before,err:=walSize(path);if err!=nil{return err}
 if before<threshold{return nil}
 if before>=walCriticalBytes{log.Printf("database WAL CRITICAL: %d MiB; checkpoint attempt",before>>20)}else{log.Printf("database WAL warning: %d MiB; checkpoint attempt",before>>20)}
 checkCtx,cancel:=context.WithTimeout(ctx,30*time.Second);defer cancel()
 busy,frames,checkpointed,err:=checkpointWAL(checkCtx,db)
 if err!=nil{return fmt.Errorf("SQLite WAL checkpoint: %w",err)}
 if busy!=0{return fmt.Errorf("SQLite WAL checkpoint busy=%d frames=%d checkpointed=%d",busy,frames,checkpointed)}
 after,err:=walSize(path);if err!=nil{return err}
 log.Printf("database WAL checkpoint OK: before=%d MiB after=%d MiB frames=%d checkpointed=%d",before>>20,after>>20,frames,checkpointed)
 return nil
}
func runDatabaseMaintenance(ctx context.Context,db *sql.DB,path string){
 clean:=func(){
  cleanCtx,cancel:=context.WithTimeout(ctx,60*time.Second);defer cancel()
  removed,err:=pruneOldSamples(cleanCtx,db,time.Now())
  if err!=nil{if ctx.Err()==nil{log.Printf("database retention failed: %v",err)};return}
  if removed>0{log.Printf("database retention removed %d expired samples",removed)}
 }
 clean()
 if n,err:=walSize(path);err==nil&&n>=walWarningBytes{
  log.Printf("database WAL at startup: %d MiB (threshold %d MiB)",n>>20,walWarningBytes>>20)
 }else if err!=nil{log.Printf("database WAL stat: %v",err)}
 checkpointTicker:=time.NewTicker(10*time.Minute)
 retentionTicker:=time.NewTicker(24*time.Hour)
 defer checkpointTicker.Stop();defer retentionTicker.Stop()
 for{
  select {
  case <-ctx.Done():return
  case <-checkpointTicker.C:
   if err:=maybeCheckpointWAL(ctx,db,path,walWarningBytes);err!=nil&&ctx.Err()==nil{log.Printf("database maintenance: %v",err)}
  case <-retentionTicker.C:clean()
  }
 }
}
