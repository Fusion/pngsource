# modern ui
!include MUI2.nsh

!define MUI_PAGE_HEADER_TEXT "PNGSource"
!define MUI_PAGE_HEADER_SUBTEXT "Embed source code in png files"
!define MUI_WELCOMEPAGE_TITLE "Embed source code in png files"
!define MUI_WELCOMEPAGE_TEXT "You are about to install PNGSource."

!define MUI_FINISHPAGE_NOAUTOCLOSE

!define MUI_FINISHPAGE_TITLE "PNGSource successfully installed."
!define MUI_FINISHPAGE_LINK "Click here to visit the project's page"
!define MUI_FINISHPAGE_LINK_LOCATION "https://github.com/fusion/pngsource"

!define MUI_PAGE_CUSTOMFUNCTION_LEAVE "OnPageWelcomeLeave"
!insertmacro MUI_PAGE_WELCOME

!define MUI_PAGE_CUSTOMFUNCTION_LEAVE "OnPageDirectoryLeave"
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!define MUI_WELCOMEPAGE_TEXT "You are about to uninstall PNGSource."
!define MUI_FINISHPAGE_TEXT "PNGSource was uninstalled.$\r$\nI hope it worked well for you."
!define MUI_FINISHPAGE_LINK "Click here to visit the project's page"
!define MUI_FINISHPAGE_LINK_LOCATION "https://github.com/fusion/pngsource"
!insertmacro MUI_UNPAGE_WELCOME
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES
!insertmacro MUI_UNPAGE_FINISH

VIProductVersion "0.1.0.0"
VIFileVersion "0.1.0.0"

Name "PNGSouce Installer"
OutFile "pngsource_installer.exe"
InstallDir "$LOCALAPPDATA\PNGSource"
 
# For removing Start Menu shortcut in Windows 7
RequestExecutionLevel user
 

Section
    SetOutPath $INSTDIR
 
    File pngsource-windows-4.0-amd64.exe
    File webview.dll
    File WebView2Loader.dll

    WriteUninstaller "$INSTDIR\uninstall.exe"
    CreateShortcut "$SMPROGRAMS\PNGSource Uninstaller.lnk" "$INSTDIR\uninstall.exe"
    CreateShortcut "$SMPROGRAMS\PNGSource.lnk" "$INSTDIR\pngsource-windows-4.0-amd64.exe"
SectionEnd
 
Section "uninstall"
    SetOutPath $INSTDIR
 
    Delete "$INSTDIR\uninstall.exe"
    Delete "$INSTDIR\pngsource-windows-4.0-amd64.exe"
    Delete "$INSTDIR\webview.exe"
    Delete "$INSTDIR\WebView2Loader.dll"
    Delete "$SMPROGRAMS\PNGSource Uninstaller.lnk"
    Delete "$SMPROGRAMS\PNGSource.lnk"
 
    RMDir $INSTDIR
SectionEnd

Function "OnPageWelcomeLeave"
  Return
FunctionEnd

Function "OnPageDirectoryLeave"
  Return
FunctionEnd
