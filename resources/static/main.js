$(document).ready(function() {

	// - - - - - - - - -
	// Table of contents
	// - - - - - - - - -

	if ($( "#body nav" ).html().trim()) {
		$( "#body nav" ).before( '<span>Contents</span>&nbsp;', '[<a id="show-hide-toc" href="#">hide</a>]' );
	}

	$("#show-hide-toc").click(function (event) {
		event.preventDefault();
		$("#body nav").toggle(250);
		if ($(this).text() === "hide") {
			$(this).text("show");
		} else {
			$(this).text("hide");
		}
	});

});
