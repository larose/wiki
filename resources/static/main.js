$(document).ready(function() {

  // - - - - - - - - -
  // Table of contents
  // - - - - - - - - -

  var html = $("#body nav").html();
  if (typeof html === "string" && html.trim()) {
    $("#body nav").before( '<span>Contents</span>&nbsp;', '[<a id="show-hide-toc" href="#">hide</a>]' );
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


  // - - - - - - - - - -
  // Keyboard shortcuts
  // - - - - - - - - - -

  [document, "#body-edit", "#message"].forEach(function(selector) {
    $(selector).bind("keydown", "alt+shift+e", function () {
      $("#edit").click();
    });

    $(selector).bind("keydown", "alt+shift+p", function () {
      $("#preview").click();
    });

    $(selector).bind("keydown", "alt+shift+d", function () {
      $("#diff").click();
    });

    $(selector).bind("keydown", "alt+shift+s", function () {
      $("#save").click();
    });
  });

});
