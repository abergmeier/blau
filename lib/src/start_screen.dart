import 'package:flutter/material.dart';
import 'package:flutter_webrtc/webrtc.dart';
import 'signaling.dart';

class StartScreen extends StatefulWidget {
  @override
  _StartState createState() => new _StartState();
}

class _StartState extends State<StartScreen> {

  MediaStream _localStream;
  final _localRenderer = new RTCVideoRenderer();
  final List<RTCVideoRenderer> _remoteRenderers = [
    RTCVideoRenderer(),
  ];
  final Set<RTCVideoRenderer> _selectedRenderers = {};
  var _signaling = new Signaling('8.8.8.8');

  @override
  Widget build(BuildContext context) {

    return Scaffold(
      appBar: AppBar(
        title: Text('BlAu'),
      ),
      body: new OrientationBuilder(
        builder: (context, orientation) {
          int crossAxisCount;
          if (orientation == Orientation.landscape) {
            crossAxisCount = 4;
          } else {
            crossAxisCount = 2;
          }
          return Column(
            children: <Widget>[
              GridView.count(
                crossAxisCount: 4,
                shrinkWrap: true,
                children: _buildSelected(),
              ),
              GridView.count(
                crossAxisCount: crossAxisCount,
                shrinkWrap: true,
                children: _buildVideos(),
              ),
            ]
          );
        }
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _startSession,
        tooltip: 'Start game',
        child: new Icon(Icons.phone),
      ),
    );
  }

  List<Widget> _buildSelected() {
    var selected = <Widget>[
      Container(
        width: MediaQuery.of(context).size.width,
        height: MediaQuery.of(context).size.height,
        child: RTCVideoView(_localRenderer),
      ),
    ];

    for (int i = 0; i != _selectedRenderers.length; i++) {
      selected.add(
        Stack(
          children: <Widget>[
            Container(
              width: MediaQuery.of(context).size.width,
              height: MediaQuery.of(context).size.height,
              child: RTCVideoView(_remoteRenderers[i]),
            ),
            Checkbox(
              value: true,
              onChanged: (bool ignore) {
                var renderer = _remoteRenderers[i];
                _unselect(renderer);
              },
            ),
          ],
        )
      );
    }
  }

  _select(RTCVideoRenderer renderer) {
    if (_selectedRenderers.add(renderer)) {
      print("Select state is wrong");
    }
  }

  _unselect(RTCVideoRenderer renderer) {
    _selectedRenderers.remove(renderer);
  }

  List<Widget> _buildVideos() {
    var videos = <Widget>[];

    for (int i = 0; i != _remoteRenderers.length; i++) {
      videos.add(
        Stack(
          children: <Widget>[
            Container(
              width: MediaQuery.of(context).size.width,
              height: MediaQuery.of(context).size.height,
              child: RTCVideoView(_remoteRenderers[i]),
            ),
            Checkbox(
              value: false,
              onChanged: (bool value) {
                var renderer = _remoteRenderers[i];
                if (value) {
                  _select(renderer);
                } else {
                  _unselect(renderer);
                }
              },
            ),
          ]
        )
      );
    }

    return videos;
  }

  @override
  initState() {
    super.initState();
    initRenderers();
  }

  initRenderers() async {
    await _localRenderer.initialize();

    _startCamera();

    _signaling..connect();
  }

  _addRemoteRenderer() async {
    var renderer = new RTCVideoRenderer();
    await renderer.initialize();
    _remoteRenderers.add(renderer);
  }

  _startCamera() async {

    final Map<String, dynamic> mediaConstraints = {
      "audio": true,
      "video": {
        "mandatory": {
          "minWidth":'320',
          "minHeight": '240',
          "minFrameRate": '15',
        },
        "facingMode": "user",
        "optional": [],
      }
    };

    try {
      navigator.getUserMedia(mediaConstraints).then((stream){
        _localStream = stream;
        _localRenderer.srcObject = _localStream;
      });
    } catch (e) {
      print(e.toString());
    }
  }

  _startSession() {
    //Navigator.of(context).push(MaterialPageRoute(builder: (context) => MyRecord("WonderWorld")));
  }
}