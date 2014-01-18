$(document).ready(function(){
    $('.search').focus(function () {
        $(this).animate({ width: "200px" }, 500);
    });

    $('.search').blur(function () {
    	$(this).animate({ width: "73px" }, 500);
    });
});