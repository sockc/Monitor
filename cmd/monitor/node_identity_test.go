package main

import (
 "context"
 "crypto/sha256"
 "encoding/hex"
 "path/filepath"
 "strings"
 "testing"
)

func TestNodeIDAutoGenerationAndDisplayNameIsolation(t *testing.T){
 ctx:=context.Background()
 db,err:=openDB(filepath.Join(t.TempDir(),"nodes.db"));if err!=nil{t.Fatal(err)};defer db.Close()
 a,tok,err:=createRandomNodeCredential(ctx,db);if err!=nil{t.Fatal(err)}
 b,tok2,err:=createRandomNodeCredential(ctx,db);if err!=nil{t.Fatal(err)}
 if a==b || !strings.HasPrefix(a,"n-")||len(a)!=22||len(tok)!=64||len(tok2)!=64 {t.Fatalf("invalid generated IDs: %s %s",a,b)}
 if !validNodeToken(ctx,db,a,tok,""){t.Fatal("generated secret not accepted")}
 if _,err=db.Exec("INSERT INTO node_metadata(name,display_name,group_name,location,notes) VALUES(?,?,?,?,?)",a,"Original","group","US","notes");err!=nil{t.Fatal(err)}
 if _,err=db.Exec("UPDATE node_metadata SET display_name=? WHERE name=?", "首页新名字",a);err!=nil{t.Fatal(err)}
 if !validNodeToken(ctx,db,a,tok,""){t.Fatal("editing display name invalidated credentials")}
 var label string
 if err=db.QueryRow("SELECT display_name FROM node_metadata WHERE name=?",a).Scan(&label);err!=nil||label!="首页新名字"{t.Fatal("display name not saved")}
 rotated,err:=rotateNodeCredential(ctx,db,a);if err!=nil{t.Fatal(err)}
 if rotated==tok||validNodeToken(ctx,db,a,tok,"")||!validNodeToken(ctx,db,a,rotated,""){t.Fatal("credential rotation broken")}
 if !validNodeToken(ctx,db,b,tok2,""){t.Fatal("other node token changed")}
 if err=db.QueryRow("SELECT display_name FROM node_metadata WHERE name=?",a).Scan(&label);err!=nil||label!="首页新名字"{t.Fatal("rebind changed metadata")}
 hash:=sha256.Sum256([]byte(rotated))
 var stored string
 if err=db.QueryRow("SELECT hash FROM agent_tokens WHERE name=?",a).Scan(&stored);err!=nil||stored!=hex.EncodeToString(hash[:]){t.Fatal("secret must be stored hashed")}
}
