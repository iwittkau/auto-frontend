package frontend

var frontendTmpl = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta data-frontend-version="{{.Version}}">
    <title>{{.BackendName}} Auto Frontend</title>
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.1.0/css/bootstrap.min.css" integrity="sha384-9gVQ4dYFwwWSjIDZnLEWnxCjeSWFphJiwGPXr1jddIhOegiu1FwO5qRGvFXOdJZ4" crossorigin="anonymous">
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
        .main .form-control {
            font-family: monospace;
            font-size: 120%;
            
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
            <div class="card" id="control-btns" style="display: none;">
                <div class="card-body">
                    <div class="form-group">
                        <button style="display: none;" class="btn btn-outline-primary" id="btn-fetch-all">Fetch all</button>
                        <button style="display: none;" class="btn btn-outline-secondary" id="btn-clear-all-o">Clear outputs</button>
                        <button style="display: none;" class="btn btn-outline-secondary" id="btn-clear-all-i">Clear inputs</button>
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
    
    var backendURL = {{.BackendUrl}};
    var $btns = [];
    var $getBtns = [];
    var $setBtns = [];
    function autorun() {
        {{range .Buttons}}
        var btn = $('#btn-{{.ID}}')
        setupButton(btn, {{.Type}}, {{.Path}}, {{.Keys}}, {{.ID}}, {{.Label}});
        $btns.push(btn)
        {{if eq .Type "get_btn"}}
        $getBtns.push(btn)
        {{end}}
        {{if eq .Type "set_btn"}}
        $setBtns.push(btn)
        {{end}}
        {{end}}
        
        if ($getBtns.length > 0) {
            $('#control-btns').show()
            $('#btn-fetch-all').show()
            $('#btn-fetch-all').click(function() {
                $.each($getBtns, function(key, value){
                    $(value).click();
                })
            })
        }
        if ($getBtns.length > 0) {
            $('#control-btns').show()
            $('#btn-clear-all-o').show()
            $('#btn-clear-all-o').click(function() {
                $.each($('[id^=output]'), function(key, value){
                    $(value).val('');
                    $(value).animate({opacity:"0.2"},100).animate({opacity:"1"},100);
                })
            })
        }
        if ($setBtns.length > 0) {
            $('#control-btns').show()
            $('#btn-clear-all-i').show()
            $('#btn-clear-all-i').click(function() {
                $.each($('[id^=input]'), function(key, value){
                    $(value).val('');
                    $(value).animate({opacity:"0.2"},100).animate({opacity:"1"},100);
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
                            $('#output-'+id+'-'+value).val(JSON.stringify(result[value]))
                            $('#output-'+id+'-'+value).animate({opacity:"0.2"},100).animate({opacity:"1"},100);
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
                        $.each( dkeys, function( key, value ) {
                            $('#input-'+id+'-'+value).animate({opacity:"0.2"},100).animate({opacity:"1"},100);
                        });
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
        $('#btn-fetch-all').prop('disabled', true);
        $('#btn-clear-all-o').prop('disabled', true);
        $('#btn-clear-all-i').prop('disabled', true);
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
