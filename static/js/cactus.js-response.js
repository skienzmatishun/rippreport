To convert the provided JavaScript code into a Zaraz component, we'll focus on organizing it into functions or classes that can be registered and executed within a Zaraz environment. Zaraz components typically involve specific lifecycle methods to manage initialization and execution when integrated into a web application. Here is an outline for that:

```javascript
// Zaraz Component
(function() {
    var _____WB$wombat$assign$function_____ = function(name) {
        return (self._wb_wombat && self._wb_wombat.local_init && self._wb_wombat.local_init(name)) || self[name];
    };

    var MyZarazComponent = {
        name: "MyZarazComponent",

        init: function(flags) {
            flags.storedSession = JSON.parse(localStorage.getItem("cactus-session"));
            var node = flags.node;
            delete flags.node;

            if (typeof node === "string") {
                node = document.querySelector(node);
            }

            this.elmApp = n.Elm.Main.init({ node: node, flags: flags });
            this.elmApp.ports.storeSession.subscribe(this.storeSessionHandler.bind(this));
        },

        storeSessionHandler: function(sessionData) {
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
```

### Explanation:
1. **Self-Executing Function**: The whole component is encapsulated in an immediately invoked function expression (IIFE) to avoid polluting the global scope.
   
2. **Component Structure**: A simple object `MyZarazComponent` is created to define the component's name and methods.

3. **Initialization (`init`)**: This method initializes the component by parsing any flags passed into it. It sets up the element where the Elm application will run.

4. **Port Subscription**: The `storeSessionHandler` method subscribes to the `storeSession` port, which listens for messages from the Elm application and saves the session data to local storage.

5. **Zaraz Registration**: At the end of the script, it checks if Zaraz is available and registers this component.

This structure is suitable for Zaraz and encapsulates the functionality you provided for future reuse and integration. Ensure the environment where this code will run supports Zaraz properly.
