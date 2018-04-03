package frontend

var frontendTmpl = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta data-frontend-version="{{.Version}}">
    <title>{{.BackendName}} Auto Frontend</title>
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-beta.3/css/bootstrap.min.css" integrity="sha384-Zug+QiDoJOrZ5t4lssLdxGhVrurbmBWopoEl+M6BdEfwnCJZtKxi1KgxUyJq13dy" crossorigin="anonymous"></link>
    <style>
        .main {
            margin-top: 1em;
            margin-bottom: 2em;
        }

        .card {
            margin-top: 1em;
        }

        .form-control:disabled {
            color: black;
        }
        .loader {
            margin-top: 0.5em;
            float: right;
            margin-left: 1em;
            border: 3px solid white;
            border-top: 3px solid #343a40 ;
            border-radius: 50%;
            width: 1.5em;
            height: 1.5em;
            animation: spin 1.5s linear infinite;
        }

        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
    </style>
</head>
<body>
<nav class="navbar navbar-dark bg-dark">
        <a class="navbar-brand" href="#">{{.BackendName}}</a> 
        <span class="navbar-text" id='connection-ok'>
            <span class="badge badge-success">Connected</span>
        </span>
        <div id='connection-reconnecting' style='display: none;'>
            <span class="navbar-text" >
                Reconnecting
           </span>
           <div class="loader"></div>
        </div>
        
</nav>

<div class="container-fluid">
<div class="main">
<div id='fetchAll' style="display: none;">
        <div class="card">
                <div class="card-body">
                    <div class="form-group">
                        <button class="btn btn-primary" id="btn-fetch-all">Fetch all</button>
                    </div>
                </div>
        </div>
</div>
{{range $btn := .Buttons}}
<div class="card">
<div class="card-body">
    {{if eq .Type "get_btn"}}
        {{range $key := .Keys}}
        <div class="form-group">
            <label for="output-{{$btn.ID}}-{{$key}}">{{$key}}</label>
            <input class="form-control" id="output-{{$btn.ID}}-{{$key}}" disabled>
        </div>
        {{end}}
    {{end}}
    
    {{if eq .Type "set_btn"}}
        {{range $key := .Keys}}
        <div class="form-group">
            <label for="input-{{$btn.ID}}-{{$key}}">{{$key}}</label>
            <input class="form-control" id="input-{{$btn.ID}}-{{$key}}">
        </div>
        {{end}}
    {{end}}
    <div class="form-group">
        <button class="btn btn-{{.Class}}" id="btn-{{.ID}}" data-method="{{.Method}}" data-path="{{.Path}}" data-keys="{{.Keys}}" data-id="{{.ID}}">{{.Label}}</button>
    </div>
</div>
</div>

{{end}}
</div>
</div>
</body>

<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.7.2/jquery.min.js" type="text/javascript"></script>
<script type="text/javascript">
    "use strict";

    $('meta').each(function(key, value){
        if ($(value).attr('data-frontend-version') !== undefined) {
            if ($(value).attr('data-frontend-version') !== {{.Version}}) {
                // location.reload(true)
                console.log("old version detected")
            }
        }
    
    });
    

    // if () {
    //     location.reload(true);
    // }
    var backendURL = {{.BackendUrl}};
    var $btns = [];
    var $getBtns = [];
    function autorun()
    {
        {{range .Buttons}}
            var btn = $('#btn-{{.ID}}')
            setupButton(btn, {{.Type}}, {{.Path}}, {{.Keys}}, {{.ID}}, {{.Label}});
            $btns.push(btn)
            {{if eq .Type "get_btn"}}
                $getBtns.push(btn)
            {{end}}
        {{end}}

        if ($getBtns.length > 0) {
            $('#fetchAll').show()
            $('#btn-fetch-all').click(function() {
                $.each($getBtns, function(key, value){
                    $(value).click();
                })
            })
        }

        setupWs()


    }

    function setupButton(btn, type, path, keys, id, label) {

        switch (type) {
            case "get_btn":
                $(btn).click(function () {
                    $.ajax({
                        type: "GET",
                        url: backendURL + path,
                        success: function (e) {
                            var dkeys = $(btn).attr('data-keys');
                            dkeys = dkeys.replace("[", "").replace("]", "");
                            dkeys = dkeys.split(" ");

                            var data ={};

                            var result = JSON.parse(e)

                            $.each( dkeys, function( key, value ) {
                                $('#output-'+id+'-'+value).val(result[value])
                            });

                            $('#output-'+ id).val(e)
                        },
                        error: function (result) {
                            console.log(result)
                        }
                    });
                });
                break;
            case "do_btn":
                $(btn).click(function () {
                    $.ajax({
                        type: "GET",
                        url: backendURL + path,
                        success: function (e) {
                            
                        },
                        error: function (result) {
                            console.log(result)
                        }
                    });
                });
                break;
            case "set_btn":
                $(btn).click(function () {
                    var dkeys = $(btn).attr('data-keys');
                    dkeys = dkeys.replace("[", "").replace("]", "");
                    dkeys = dkeys.split(" ");

                    var data ={};

                    $.each( dkeys, function( key, value ) {
                        data[value] = $('#input-'+id+'-'+value).val()
                    });

                    $.ajax({
                        type: "POST",
                        url: backendURL + path,
                        success: function (e) {
                            $('#output-'+ id).val(e)
                        },
                        error: function (result) {
                            console.log(result)
                        },
                        data: JSON.stringify(data)
                    });
                });
                break;

        }
    }

    function setupWs() {
        var wsUrl = window.location.href
        wsUrl = wsUrl.replace("http://", "ws://")
        var ws = new WebSocket(wsUrl+'ws')

        ws.onerror = function (error) {
            reconnectWs(wsUrl); 
        };

        ws.onclose = function (error) {
            reconnectWs(wsUrl);            
        };

        
    }

    function reconnectWs(wsUrl) {
        $.each($btns, function(key, value){
                $(value).prop('disabled', true);
            });
            $('#connection-reconnecting').show()
            $('#connection-ok').hide()
            setInterval(function() {
                var reconnect = new WebSocket(wsUrl+'ws')
                reconnect.onopen = function() {
                    location.reload();
                }
            }, 1000);
    }

    if (document.addEventListener) document.addEventListener("DOMContentLoaded", autorun, false);
    else if (document.attachEvent) document.attachEvent("onreadystatechange", autorun);
    else window.onload = autorun;
</script>
</html>`
