@echo off

:init
set project_root=%~dp0
cd /d %project_root%

:select
goto :protoc

goto :end

:protoc
protoc --go_out=../../proto-message ../*.proto -I../
goto :end

:end
cd /d %project_root%
goto :eof