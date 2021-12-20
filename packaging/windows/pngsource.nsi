# modern ui
!include MUI2.nsh

Name "PNGSouce Installer"

# define name of installer
OutFile "pngsource_installer.exe"
 
# define installation directory
InstallDir "$LOCALAPPDATA\PNGSource"
 
# For removing Start Menu shortcut in Windows 7
RequestExecutionLevel user
 
# start default section
Section
 
    # set the installation directory as the destination for the following actions
    SetOutPath $INSTDIR
 
    # Copy file
    File pngsource.exe
    File webview.dll
    File WebView2Loader.dll

    # create the uninstaller
    WriteUninstaller "$INSTDIR\uninstall.exe"
 
    # create a shortcut named "new shortcut" in the start menu programs directory
    # point the new shortcut at the program uninstaller
    CreateShortcut "$SMPROGRAMS\uninstall pngsource.lnk" "$INSTDIR\uninstall.exe"
SectionEnd
 
# uninstaller section start
Section "uninstall"
 
    # first, delete the uninstaller
    Delete "$INSTDIR\uninstall.exe"
 
    # second, remove the link from the start menu
    Delete "$SMPROGRAMS\uninstall pngsource.lnk"
 
    RMDir $INSTDIR
# uninstaller section end
SectionEnd
