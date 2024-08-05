function reloadCss()
{
    document.querySelectorAll("link[rel=stylesheet]").forEach(link => link.href = link.href.replace(/\?.*|$/, "?" + Date.now()))
}

function connect() {
	var ws = new WebSocket("ws://localhost:3000/__rawdog-md/watcher");
	var autoReconnect = true;

	// automatically attempt to reconnect on close
	ws.onclose = function () {
		if (autoReconnect) {
			setTimeout(function () {
				connect();
			}, 1000);
		}
	};

	ws.onmessage = function (event) {
		if (event.data === "reload") {
			newPage = fetch(window.location.href)
				.then(response => response.text())
				.then(html => {
					document.open();
					document.write(html);
					document.close();
					// TODO: separate css reload from page reload
					//reloadCss();
					autoReconnect = false;
					ws.close();
				});
		} 
	};
}

connect();