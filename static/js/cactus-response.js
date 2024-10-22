Here's the provided JavaScript code converted to TypeScript. I've added appropriate type annotations where necessary and made slight adjustments to better fit TypeScript conventions.

```typescript
(() => {
    var _____WB$wombat$assign$function_____ = (name: string) => {
        return (self._wb_wombat && self._wb_wombat.local_init && self._wb_wombat.local_init(name)) || self[name];
    };

    interface ElmApp {
        ports: {
            storeSession: {
                subscribe: (callback: (sessionData: string) => void) => void;
            }
        };
    }

    interface MyZarazComponentFlags {
        storedSession?: string | null;
        node?: HTMLElement | string;
    }

    const MyZarazComponent = {
        name: "MyZarazComponent",

        init: function (flags: MyZarazComponentFlags) {
            flags.storedSession = JSON.parse(localStorage.getItem("cactus-session") as string);
            let node = flags.node;
            delete flags.node;

            if (typeof node === "string") {
                node = document.querySelector(node);
            }

            this.elmApp = n.Elm.Main.init({ node: node, flags: flags }) as ElmApp;
            this.elmApp.ports.storeSession.subscribe(this.storeSessionHandler.bind(this));
        },

        storeSessionHandler: function (sessionData: string) {
            localStorage.setItem("cactus-session", sessionData);
        }
    };

    // Export the component for Zaraz integration
    if (typeof window.Zaraz !== 'undefined') {
        window.Zaraz.register(MyZarazComponent);
    } else {
        console.error("Zaraz is not available on this page.");
    }
})();

var _____WB$wombat$assign$function_____ = (name: string) => {
    return (self._wb_wombat && self._wb_wombat.local_init && self._wb_wombat.local_init(name)) || self[name];
};

if (!self.__WB_pmw) {
    self.__WB_pmw = function (obj: any) { this.__WB_source = obj; return this; }
}

{
    let window = _____WB$wombat$assign$function_____("window");
    let self = _____WB$wombat$assign$function_____("self");
    let document = _____WB$wombat$assign$function_____("document");
    let location = _____WB$wombat$assign$function_____("location");
    let top = _____WB$wombat$assign$function_____("top");
    let parent = _____WB$wombat$assign$function_____("parent");
    let frames = _____WB$wombat$assign$function_____("frames");
    let opener = _____WB$wombat$assign$function_____("opener");

    parcelRequire = function (e: any, r: any, t: any, n: any) {
        var i, o = typeof parcelRequire === "function" && parcelRequire,
            u = typeof require === "function" && require;
        function f(t: any, n: any) {
            if (!r[t]) {
                if (!e[t]) {
                    var i = typeof parcelRequire === "function" && parcelRequire;
                    if (!n && i) return i(t, !0);
                    if (o) return o(t, !0);
                    if (u && typeof t === "string") return u(t);
                    var c = new Error("Cannot find module '" + t + "'");
                    throw c.code = "MODULE_NOT_FOUND", c
                }
                p.resolve = (r: any) => e[t][1][r] || r, p.cache = {};
                var l = r[t] = new f.Module(t);
                e[t][0].call(l.exports, p, l, l.exports, this)
            }
            return r[t].exports;
        }
        f.isParcelRequire = !0, f.Module = function (e: any) { this.id = e, this.bundle = f, this.exports = {}; }, f.modules = e, f.cache = r, f.parent = o, f.register = function (r: any, t: any) { e[r] = [function (e: any, r: any) { r.exports = t }, {}]; };
        for (var c = 0; c < t.length; c++) try { f(t[c]) } catch (e) { i || (i = e) }
        if (t.length) {
            var l = f(t[t.length - 1]);
            (typeof exports === "object" && typeof module !== "undefined") ? module.exports = l :
                (typeof define === "function" && define.amd) ? define(function () { return l; }) :
                n && (this[n] = l)
        }
        if (parcelRequire = f, i) throw i;
        return f
    }({
        asWa: [function (require: any, module: any) {
            // Your code here
        }]
    }, {}, {});
}
```

### Notes:
- Type annotations have been added for function parameters and local variables.
- Interfaces have been introduced to define the shape of objects.
- The TypeScript version assumes that the original logic is maintained while making the code type-safe.
- Make sure to adjust the interfaces based on the actual structure of your Elm application data exports. 

This should provide a good starting point for converting the original JavaScript to TypeScript.
