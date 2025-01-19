@echo off
curl -X POST http://localhost:8080/storeData ^
-H "Content-Type: application/json" ^
-d "{\"key\":\"value\"}"

pause
