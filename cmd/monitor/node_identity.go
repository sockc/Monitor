package main

import (
 "context"
 "crypto/rand"
 "crypto/sha256"
 "database/sql"
 "encoding/hex"
 "errors"
 "fmt"
)

// New nodes are identified by random, immutable IDs. User-facing labels are
// stored separately in node_metadata and never participate in authentication.
func createRandomNodeCredential(ctx context.Context, db *sql.DB)(string,string,error){
 for attempt:=0;attempt<8;attempt++{
  idBytes:=make([]byte,10)
  secretBytes:=make([]byte,32)
  if _,e:=rand.Read(idBytes);e!=nil{return "","",e}
  if _,e:=rand.Read(secretBytes);e!=nil{return "","",e}
  name:="n-"+hex.EncodeToString(idBytes)
  secret:=hex.EncodeToString(secretBytes)
  hash:=sha256.Sum256([]byte(secret))
  result,e:=db.ExecContext(ctx,"INSERT OR IGNORE INTO agent_tokens(name,hash) VALUES(?,?)",name,hex.EncodeToString(hash[:]))
  if e!=nil{return "","",e}
  added,e:=result.RowsAffected()
  if e!=nil{return "","",e}
  if added==1{return name,secret,nil}
 }
 return "","",fmt.Errorf("failed to allocate unique node ID")
}

// Rebinding is separate from editing metadata. It only rotates the secret for
// the selected immutable node identity, preserving history, location and plan.
func rotateNodeCredential(ctx context.Context,db *sql.DB,name string)(string,error){
 if !validNodeName(name){return "",fmt.Errorf("invalid node identity")}
 var existing string
 if e:=db.QueryRowContext(ctx,"SELECT hash FROM agent_tokens WHERE name=?",name).Scan(&existing);e!=nil&& !errors.Is(e,sql.ErrNoRows){return "",e}
 secretBytes:=make([]byte,32)
 if _,e:=rand.Read(secretBytes);e!=nil{return "",e}
 secret:=hex.EncodeToString(secretBytes)
 hash:=sha256.Sum256([]byte(secret))
 _,e:=db.ExecContext(ctx,"INSERT INTO agent_tokens(name,hash) VALUES(?,?) ON CONFLICT(name) DO UPDATE SET hash=excluded.hash",name,hex.EncodeToString(hash[:]))
 if e!=nil{return "",e}
 return secret,nil
}
