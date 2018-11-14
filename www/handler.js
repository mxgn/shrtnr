$(function () {
    $(".button").click(function () {
        // validate and process form here

        $('.error').hide();
        var name = $("input#long-url").val();
        if (name == "") {
            $("label#url_error").show();
            $("input#url").focus();
            return false;
        }

    });
});