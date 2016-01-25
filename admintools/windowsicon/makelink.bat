IF EXIST "%userprofile%\accordsw" (GOTO _CopyIcon)
mkdir "%userprofile%\accordsw"

:_CopyIcon
IF EXIST "%userprofile%\accordsw\accordpb.ico" (GOTO _CopyScript)
copy accordpb.ico "%userprofile%\accordsw\"

:_CopyScript
IF EXIST "%userprofile%\accordsw\makelink.bat" (GOTO _MakeLink)
copy makelink.bat "%userprofile%\accordsw\"

:_MakeLink
(echo [InternetShortcut]
echo URL=http://ec2-52-23-176-52.compute-1.amazonaws.com:8250/
echo IconFile="%userprofile%\accordsw\accordpb.ico"
echo IconIndex=0) >"%userprofile%\desktop\Accord Office.url"