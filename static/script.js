const themeSwitch = document.getElementById("themeSwitch");
const themeLabel = document.getElementById("themeLabel");
function setTheme(dark) {
	if (dark) {
		document.body.classList.add("dark");
	} else {
		document.body.classList.remove("dark");
	}
}
themeSwitch.addEventListener("change", (e) => {
	setTheme(e.target.checked);
	localStorage.setItem("uism_theme", e.target.checked ? "dark" : "light");
});
(function () {
	const saved = localStorage.getItem("uism_theme");
	const prefersDark = window.matchMedia(
		"(prefers-color-scheme: dark)"
	).matches;
	if (saved === "dark" || (!saved && prefersDark)) {
		themeSwitch.checked = true;
		setTheme(true);
	}
})();
document.addEventListener("DOMContentLoaded", function () {
	const copyBtn = document.querySelector(".copy-btn");
	const qrBtn = document.querySelector(".qr-btn");
	const shortUrl = document.getElementById("shortUrlLink");
	const qrModal = document.getElementById("qr-modal");
	const qrClose = document.querySelector(".qr-close");
	const qrCodeDiv = document.getElementById("qr-code");

	if (copyBtn && shortUrl) {
		copyBtn.onclick = function () {
			navigator.clipboard.writeText(shortUrl.href);
			copyBtn.innerHTML = "✓";
			setTimeout(() => {
				copyBtn.innerHTML = `<svg viewBox="0 0 20 20"><rect x="7" y="3" width="9" height="14" rx="2"/><rect x="3" y="7" width="9" height="10" rx="2"/></svg>`;
			}, 1200);
		};
	}
	if (qrBtn && shortUrl) {
		qrBtn.onclick = function () {
			qrCodeDiv.innerHTML = "";
			const accent =
				getComputedStyle(document.body).getPropertyValue("--accent") ||
				"#3b82f6";
			const qr = new QRious({
				element: document.createElement("canvas"),
				value: shortUrl.href,
				size: 180,
				background: "transparent",
				foreground: accent.trim(),
			});
			qrCodeDiv.appendChild(qr.element);
			qrModal.style.display = "flex";
		};
	}
	if (qrClose) {
		qrClose.onclick = function () {
			qrModal.style.display = "none";
		};
	}
	qrModal &&
		qrModal.addEventListener("click", function (e) {
			if (e.target === qrModal) qrModal.style.display = "none";
		});
});
