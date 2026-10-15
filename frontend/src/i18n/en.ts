export const en: Record<string, string> = {
  // Title bar
  'title.minimize': 'Minimize',
  'title.maximize': 'Maximize / Restore',
  'title.close': 'Close',

  // Top bar
  'sdk.label': 'Android SDK',
  'sdk.notSet': 'Not set',
  'sdk.pick': 'Choose folder',
  'sdk.manual': 'Edit path',
  'sdk.placeholder': 'e.g. C:\\Users\\<you>\\AppData\\Local\\Android\\Sdk',
  'sdk.apply': 'Apply',
  'sdk.ready': '✓ adb ready ({version})',
  'sdk.notReady': '✕ adb not ready',
  'sdk.toggleTheme': 'Toggle theme',

  // Sidebar
  'sidebar.title': 'Devices',
  'sidebar.create': 'New phone profile',
  'sidebar.delete': 'Delete this profile',
  'sidebar.serialAuto': 'serial empty · auto-match',
  'sidebar.empty': 'No phone profiles',
  'sidebar.emptyHint': 'Click + in the corner to create one',
  'sidebar.connUnknown': 'Connection unknown — use "Check connection"',

  // Device panel
  'devices.title': 'adb devices',
  'devices.refreshing': 'Refreshing…',
  'devices.refresh': 'Refresh',
  'state.device': 'online·authorized',
  'state.unauthorized': 'unauthorized',
  'state.offline': 'offline',
  'state.recovery': 'Recovery',
  'state.sideload': 'Sideload',
  'devices.none': 'No devices found — connect a phone and click "Refresh"',
  'devices.adbNotReady': 'adb not ready — set the Android SDK directory first',

  // Form
  'form.basic': 'Basic info',
  'form.package': 'APK package',
  'form.options': 'Install options',
  'form.name': 'Phone name',
  'form.nameHint': '(= yml file name)',
  'form.namePlaceholder': 'e.g. Office test phone',
  'form.serial': 'Device serial',
  'form.serialHint': '(first column of adb devices; empty = auto-match the only online device)',
  'form.serialPlaceholder': 'e.g. emulator-5554 / 1A2B3C4D',
  'form.apkPath': 'APK path',
  'form.apkPathHint': '(absolute path, .apk)',
  'form.apkPlaceholder': 'D:\\apk\\demo-release.apk',
  'form.browse': 'Browse…',
  'form.reinstall': 'Reinstall',
  'form.reinstallHint': '-r keep data on upgrade',
  'form.downgrade': 'Allow downgrade',
  'form.grantPerms': 'Grant runtime permissions',
  'form.autoLaunch': 'Auto-launch after install',
  'form.launchPackagePlaceholder': 'package name, e.g. com.example.app',
  'form.notesPlaceholder': 'Notes (optional)',
  'form.hint': 'Two pre-checks run before install: ① adb availability → ② device reachable; install is skipped if either fails.',
  'form.empty': 'Select or create a phone on the left to start',
  'check.ok': '✓ Connected: ',
  'check.fail': '✕ Failed: ',

  // Actions
  'action.save': 'Save',
  'action.check': 'Check connection',
  'action.checking': 'Checking…',
  'action.install': 'Check & Install',
  'action.installing': 'Installing…',

  // Logs
  'log.title': 'Logs',
  'log.openConfigs': 'Open configs folder',
  'log.clear': 'Clear logs',
  'log.empty': 'No logs yet',

  // Modals
  'common.cancel': 'Cancel',
  'common.confirm': 'OK',
  'modal.createTitle': 'New phone profile',
  'modal.createBody': 'Will create "{name}.yml" under configs/ — one file per phone.',
  'modal.namePlaceholder': 'Phone name, e.g. Office test phone',
  'modal.deleteTitle': 'Delete phone profile',
  'modal.deleteBody': 'Delete "{name}"?',
  'modal.deleteBodyHint': 'This removes configs/{name}.yml and cannot be undone.',

  // Toast / log messages
  'toast.saved': 'Saved configs/{name}.yml',
  'toast.created': 'Created configs/{name}.yml',
  'toast.deleted': 'Deleted configs/{name}.yml',
  'toast.needSdkDir': 'Enter the SDK directory path first',
  'toast.needPhoneName': 'Enter a phone name',
  'log.devicesRefreshed': 'Device list refreshed: {n} device(s)',
  'log.refreshFailed': 'Failed to refresh devices: {msg}',
  'log.initDone': 'Initialized: {n} phone profile(s); adb ready ({version})',
  'log.initDoneNotReady': 'Initialized: {n} phone profile(s); adb not ready — set the Android SDK directory in the top bar',
  'log.initFailed': 'Initialization failed: {msg}',
}
