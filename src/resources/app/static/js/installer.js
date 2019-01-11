/*
  Handles the installer interface
 */
const remote = require('electron').remote;
let installer = {
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
        // astilectron.sendMessage({
        //   name: "ready",
        //   payload: ""
        // }, function(message) {
        //
        // });

        installer.bindEvents();
        installer.listen();
      })
  },
  listen: function() {
    astilectron.onMessage(function(message) {
      var parsed = $.parseJSON(message.payload);
      switch (message.name) {
        case "state":

          // Install state

          // TODO: If state == 'confirm-av'
          // TODO: Show antivirus warning to exclude
          //
          // TODO: If state == 'done'
          // TODO: Show done
          //
          // TODO: If state == 'error'
          // TODO: Show error

          break;

        case "update":

          // Install update?

          break;

        default:
          console.log("Unknown command '" + message.name + "' received");
          break;
        }
      });
  },
  // install the miner services
  install: function(rigName, installPath) {
    if (rigName == '')
    {
      alert('No rig name set');
      return false;
    }

    if (installPath == '')
    {
      alert('No install path set');
      return false;
    }

    astilectron.sendMessage({
      name: "install",
      payload: {
        rigName: rigName,
        installPath: installPath,
      },
    }, function(message) {
      var data = message.payload;
      if (data.status == 'error')
      {
        alert('Unable to install: ' + data.message);
      }
      else if (data.status == 'ok')
      {
        // TODO Show the installation progress
        console.log('show install progress');
      }
    });

  },
  // Bind to UI events using jQuery
  bindEvents: function() {
    var totalSteps = 5;
    var currentStep = 1;

    $('.wizard-continue').bind('click', function(){
      var buttonRole = $(this).data('role');
      var buttonStep = $(this).data('step');
      var nextStep = buttonStep + 1;
      var previousStep = buttonStep - 1;

      if (buttonRole == 'next')
      {
        if (nextStep == 3) // Next from rig name page
        {
          if ($('#rig_name').val() == '')
          {
            $('#rig_name_error').removeClass('d-none');
            return false;
          }
          else
          {
            $('#rig_name_error').addClass('d-none');
          }
        }

        if (nextStep == 4) // Next from install path selection
        {
          if ($('#install_path').val() == '')
          {
            $('#install_path_error').removeClass('d-none');
            return false;
          }
          else
          {
            $('#install_path_error').addClass('d-none');
          }
        }

        if (nextStep == 5) // Confirmed
        {
          var installing = installer.install($('#rig_name').val(), $('#install_path').val());
          if (installing == false) {
            return false;
          }
        }

        if (nextStep < totalSteps)
        {
          $('#step_' + buttonStep).toggleClass('hide');
          $('#step_' + nextStep).toggleClass('hide');
          $('#step_' + buttonStep + '_index').toggleClass('text-white');
          $('#step_' + nextStep + '_index').toggleClass('text-white');
          currentStep++;
        }
      }
      else
      {
        if (previousStep > 0)
        {
          $('#step_' + buttonStep).toggleClass('hide');
          $('#step_' + previousStep).toggleClass('hide');
          $('#step_' + buttonStep + '_index').toggleClass('text-white');
          $('#step_' + previousStep + '_index').toggleClass('text-white');
          currentStep--;
        }
      }

      if (currentStep == 4) // Confirmed page, set name and install path
      {
        $('#confirm_rig_name').html($('#rig_name').val());
        $('#confirm_install_path').html($('#install_path').val());
      }


    });

    $('#install_path_selector').bind('click', function(){
      astilectron.showOpenDialog({properties: ['openDirectory',], title: "Select your installation directory"}, function(path) {
          $('#install_path').val(path);
          $('#install_path_error').addClass('d-none');
      });
    });


  },
};
