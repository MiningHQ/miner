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
          // Show modal with error
          $('#error_list').html(message.payload.message);
          $('#error_modal').modal();
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
          if (parsed.Stats.Hashrate != undefined)
          {
            $('#current_hashrate').html(parsed.Stats.Hashrate + " H/s");
          } else $('#current_hashrate').html("0 H/s");
          $('#shares_total').html(parsed.Stats.TotalShares);
          $('#shares_accepted').html(parsed.Stats.AcceptedShares);
          $('#shares_rejected').html(parsed.Stats.RejectedShares);
          $('#rig_logs').html(parsed.HTMLLogs);

          if (parsed.State == 2) // Mining = 2;
          {
            $('#state_info').addClass('text-success');
            $('#state_info').removeClass('text-danger');
            $('#state_info').removeClass('text-warning');
            $('#state_info').html('Mining');
            $('#state_info').show();

            $('#pause').removeClass('d-none');
            $('#resume').addClass('d-none');
          }
          else if (parsed.State == 3) // StopMining = 3;
          {
            $('#state_info').removeClass('text-success');
            $('#state_info').removeClass('text-warning');
            $('#state_info').addClass('text-danger');
            $('#state_info').html('Not mining');
            $('#state_info').show();
          }
          else if (parsed.State == 4) // PauseMining = 4;
          {
            $('#state_info').removeClass('text-success');
            $('#state_info').removeClass('text-danger');
            $('#state_info').addClass('text-warning');
            $('#state_info').html('Paused');
            $('#state_info').show();
            $('#pause').addClass('d-none');
            $('#resume').removeClass('d-none');
          }

          break;

        default:
          console.log("Unknown command '" + message.name + "' received");
          break;
        }
      });
  },
  // Bind to UI events using jQuery
  bindEvents: function() {

    $('.exit').bind('click', function(){
       remote.getCurrentWindow().close();
    });

    $('.minimize').bind('click', function(){
      remote.getCurrentWindow().minimize();
    });

    $('#pause').bind('click', function(){

      astilectron.sendMessage({
        name: "pause",
        payload: ""
      }, function(message) {
        if (message.payload.status == 'error')
        {
          $('#error_list').html(message.payload.message);
          $('#error_modal').modal();
        }
        else
        {
          $('#state_info').removeClass('text-success');
          $('#state_info').removeClass('text-danger');
          $('#state_info').addClass('text-warning');
          $('#state_info').html('Paused');
          $('#state_info').show();

          // Hide pause, show resume
          $('#pause').addClass('d-none');
          $('#resume').removeClass('d-none');
        }
      });
    });

    $('#resume').bind('click', function(){

      astilectron.sendMessage({
        name: "resume",
        payload: ""
      }, function(message) {
        if (message.payload.status == 'error')
        {
          $('#error_list').html(message.payload.message);
          $('#error_modal').modal();
        }
        else
        {
          $('#state_info').addClass('text-success');
          $('#state_info').removeClass('text-danger');
          $('#state_info').removeClass('text-warning');
          $('#state_info').html('Mining');
          $('#state_info').show();

          // Hide pause, show resume
          $('#pause').removeClass('d-none');
          $('#resume').addClass('d-none');
        }
      });
    });

    $('#refresh').bind('click', function(){
      astilectron.sendMessage({name: "refresh", payload: ""}, function(message){
      });
    });
  },
};
