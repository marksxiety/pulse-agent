export namespace main {
	
	export class Bootstrap {
	    theme: string;
	    themes: theme.Theme[];
	    version: string;
	    alwaysOnTop: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Bootstrap(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.theme = source["theme"];
	        this.themes = this.convertValues(source["themes"], theme.Theme);
	        this.version = source["version"];
	        this.alwaysOnTop = source["alwaysOnTop"];
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

}

export namespace theme {
	
	export class Theme {
	    name: string;
	    bg: string;
	    surface: string;
	    border: string;
	    muted: string;
	    subtle: string;
	    text: string;
	    dim: string;
	    cpu_accent: string;
	    mem_accent: string;
	    disk_accent: string;
	    warn: string;
	    danger: string;
	
	    static createFrom(source: any = {}) {
	        return new Theme(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.bg = source["bg"];
	        this.surface = source["surface"];
	        this.border = source["border"];
	        this.muted = source["muted"];
	        this.subtle = source["subtle"];
	        this.text = source["text"];
	        this.dim = source["dim"];
	        this.cpu_accent = source["cpu_accent"];
	        this.mem_accent = source["mem_accent"];
	        this.disk_accent = source["disk_accent"];
	        this.warn = source["warn"];
	        this.danger = source["danger"];
	    }
	}

}

