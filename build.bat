@echo off
echo Building GTD Android App...

if exist "GTD Organizer.apk" del "GTD Organizer.apk"
if exist "GTD Organizer" rmdir /s /q "GTD Organizer"

echo Tidying modules...
go mod tidy

echo Building APK...
fyne package -os android -appID com.gtdandroid.app -name "GTD Organizer" -icon icon.png

if exist "GTD Organizer.apk" (
    echo Installing APK...
    adb install -r "GTD Organizer.apk"
    echo Done! App installed successfully.
) else (
    echo Build failed!
)