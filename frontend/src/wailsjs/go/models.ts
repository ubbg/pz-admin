export namespace main {
	
	export class Config {
	    theme?: string;
	    colorScheme?: string;
	    useSystemTitleBar?: boolean;
	    enableLogging?: boolean;
	    enableTrace?: boolean;
	    enableDebug?: boolean;
	    enableInfo?: boolean;
	    enableWarn?: boolean;
	    enableError?: boolean;
	    enableFatal?: boolean;
	    maxLogFiles?: number;
	    language?: string;
	    saveWindowStatus?: boolean;
	    windowStartState?: number;
	    windowStartPositionX?: number;
	    windowStartPositionY?: number;
	    windowStartSizeX?: number;
	    windowStartSizeY?: number;
	    windowScale?: number;
	    opacity?: number;
	    windowEffect?: number;
	    checkForUpdates?: boolean;
	    lastUpdateCheck?: number;
	    rememberCredentials?: boolean;
	    autoConnect?: boolean;
	    rconCheckInterval?: number;
	    disableWeatherControlButtons?: boolean;
	    disableRandomButtons?: boolean;
	    disableOtherButtons?: boolean;
	    debugMode?: boolean;
	    sandboxAccessMode?: string;
	    sandboxPath?: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.theme = source["theme"];
	        this.colorScheme = source["colorScheme"];
	        this.useSystemTitleBar = source["useSystemTitleBar"];
	        this.enableLogging = source["enableLogging"];
	        this.enableTrace = source["enableTrace"];
	        this.enableDebug = source["enableDebug"];
	        this.enableInfo = source["enableInfo"];
	        this.enableWarn = source["enableWarn"];
	        this.enableError = source["enableError"];
	        this.enableFatal = source["enableFatal"];
	        this.maxLogFiles = source["maxLogFiles"];
	        this.language = source["language"];
	        this.saveWindowStatus = source["saveWindowStatus"];
	        this.windowStartState = source["windowStartState"];
	        this.windowStartPositionX = source["windowStartPositionX"];
	        this.windowStartPositionY = source["windowStartPositionY"];
	        this.windowStartSizeX = source["windowStartSizeX"];
	        this.windowStartSizeY = source["windowStartSizeY"];
	        this.windowScale = source["windowScale"];
	        this.opacity = source["opacity"];
	        this.windowEffect = source["windowEffect"];
	        this.checkForUpdates = source["checkForUpdates"];
	        this.lastUpdateCheck = source["lastUpdateCheck"];
	        this.rememberCredentials = source["rememberCredentials"];
	        this.autoConnect = source["autoConnect"];
	        this.rconCheckInterval = source["rconCheckInterval"];
	        this.disableWeatherControlButtons = source["disableWeatherControlButtons"];
	        this.disableRandomButtons = source["disableRandomButtons"];
	        this.disableOtherButtons = source["disableOtherButtons"];
	        this.debugMode = source["debugMode"];
	        this.sandboxAccessMode = source["sandboxAccessMode"];
	        this.sandboxPath = source["sandboxPath"];
	    }
	}
	export class Coordinates {
	    x: number;
	    y: number;
	    z: number;
	
	    static createFrom(source: any = {}) {
	        return new Coordinates(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.x = source["x"];
	        this.y = source["y"];
	        this.z = source["z"];
	    }
	}
	export class Credentials {
	    ip: string;
	    port: string;
	    password: string;
	
	    static createFrom(source: any = {}) {
	        return new Credentials(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ip = source["ip"];
	        this.port = source["port"];
	        this.password = source["password"];
	    }
	}
	export class ImportOptionsResponse {
	    options: Record<string, string>;
	    success: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ImportOptionsResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.options = source["options"];
	        this.success = source["success"];
	    }
	}
	export class ItemRecord {
	    itemId: string;
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new ItemRecord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.itemId = source["itemId"];
	        this.count = source["count"];
	    }
	}
	export class Notification {
	    title: string;
	    message: string;
	    path: string;
	    variant: string;
	    parameters?: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new Notification(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.message = source["message"];
	        this.path = source["path"];
	        this.variant = source["variant"];
	        this.parameters = source["parameters"];
	    }
	}
	export class Option {
	    name: string;
	    value: string;
	    kind: string;
	
	    static createFrom(source: any = {}) {
	        return new Option(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.value = source["value"];
	        this.kind = source["kind"];
	    }
	}
	export class Player {
	    name: string;
	    online: boolean;
	    accessLevel: string;
	    banned: boolean;
	    godmode: boolean;
	    invisible: boolean;
	    noclip: boolean;
	    voiceBanned: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Player(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.online = source["online"];
	        this.accessLevel = source["accessLevel"];
	        this.banned = source["banned"];
	        this.godmode = source["godmode"];
	        this.invisible = source["invisible"];
	        this.noclip = source["noclip"];
	        this.voiceBanned = source["voiceBanned"];
	    }
	}
	export class RconResponse {
	    response: string;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new RconResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.response = source["response"];
	        this.error = source["error"];
	    }
	}
	export class SandboxAccess {
	    configured: boolean;
	    mode: string;
	    path: string;
	    file: string;
	    reason: string;
	
	    static createFrom(source: any = {}) {
	        return new SandboxAccess(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.configured = source["configured"];
	        this.mode = source["mode"];
	        this.path = source["path"];
	        this.file = source["file"];
	        this.reason = source["reason"];
	    }
	}
	export class SandboxVar {
	    key: string;
	    group: string;
	    name: string;
	    value: string;
	    kind: string;
	    quoted: boolean;
	
	    static createFrom(source: any = {}) {
	        return new SandboxVar(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.group = source["group"];
	        this.name = source["name"];
	        this.value = source["value"];
	        this.kind = source["kind"];
	        this.quoted = source["quoted"];
	    }
	}
	export class SandboxDocument {
	    success: boolean;
	    file: string;
	    vars: SandboxVar[];
	    checksum: string;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new SandboxDocument(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.file = source["file"];
	        this.vars = this.convertValues(source["vars"], SandboxVar);
	        this.checksum = source["checksum"];
	        this.error = source["error"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SandboxSaveResult {
	    success: boolean;
	    backup: string;
	    changed: string[];
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new SandboxSaveResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.backup = source["backup"];
	        this.changed = source["changed"];
	        this.error = source["error"];
	    }
	}
	
	export class ServerMessage {
	    message: string;
	    lineColors: Record<number, string>;
	
	    static createFrom(source: any = {}) {
	        return new ServerMessage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.message = source["message"];
	        this.lineColors = source["lineColors"];
	    }
	}
	export class UpdateInfo {
	    updateAvailable: boolean;
	    currentVersion: string;
	    latestVersion: string;
	    name: string;
	    releaseNotes: string;
	    downloadUrl: string;
	    releaseUrl: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.updateAvailable = source["updateAvailable"];
	        this.currentVersion = source["currentVersion"];
	        this.latestVersion = source["latestVersion"];
	        this.name = source["name"];
	        this.releaseNotes = source["releaseNotes"];
	        this.downloadUrl = source["downloadUrl"];
	        this.releaseUrl = source["releaseUrl"];
	    }
	}
	export class UpdateSource {
	    configured: boolean;
	    owner: string;
	    repo: string;
	    url: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateSource(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.configured = source["configured"];
	        this.owner = source["owner"];
	        this.repo = source["repo"];
	        this.url = source["url"];
	    }
	}

}

