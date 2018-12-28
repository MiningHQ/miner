/*
  Handles the manager interface
 */
const remote = require('electron').remote;
let manager = {
  init: function() {
      asticode.loader.init();

      shared.bindExternalLinks();

      // if (shared.isMac()) {
      //   $('.minimize').hide();
      //   $('.exit').hide();
      //   $('.settings').css('right', '15px');
      //   $('.help').css('right', '45px');
      // }

      // Wait for the ready signal
      document.addEventListener('astilectron-ready', function() {
        // Send a notification when we are ready
        astilectron.sendMessage({
          name: "ready",
          payload: ""
        }, function(message) {

        });

        manager.bindEvents();
        manager.listen();
      })
  },
  listen: function() {
    var errorCount = 0;
    astilectron.onMessage(function(message) {
      var parsed = $.parseJSON(message.payload);
      switch (message.name) {
        case "setup":
          $('#rig_name').html(parsed.name);
          $('#rig_link').attr('href', parsed.link);
          break;

        case "update":
          console.log("received update packet")
          break;

        default:
          console.log("Unknown command '" + message.name + "' received");
          break;
        }
      });
  },
  // Bind to UI events using jQuery
  bindEvents: function() {

    // $('.header-button.help').bind('click', function(){
    //   $('#help').toggleClass('dn');
    // });
    // $('.header-button.minimize').bind('click', function(){
    //   remote.getCurrentWindow().minimize();
    // });
    // $('.header-button.exit').bind('click', function(){
    //   remote.getCurrentWindow().close();
    // });

    // $('.close-help').bind('click', function(){
    //   $('#help').toggleClass('dn');
    // });

    $('#refresh').bind('click', function(){
      astilectron.sendMessage({name: "refresh", payload: configData}, function(message){

      });

    });
  },
};
