import 'package:flutter/material.dart';
import 'package:flutter_webrtc/webrtc.dart';

class Session extends StatefulWidget {

  @override
  _SessionState createState() => new _SessionState();
}

class _SessionState extends State<Session> {
  MediaStream _localStream;
  final _localRenderer = new RTCVideoRenderer();
  
  @override
  initState() {
    super.initState();
    initRenderers();
  }

  initRenderers() async {
    await _localRenderer.initialize();
  }

  void addParking(Row row, int rowI) {

    for (int i = 0; i != rowI; i++) {
      row.children.add(Container(color: Colors.grey));
    }
  }

  void addStored(Row row, int rowI) {
    row.children.add(Container(color: Colors.blue));
    row.children.add(Container(color: Colors.yellow));
    row.children.add(Container(color: Colors.red));
    row.children.add(Container(color: Colors.black));
    row.children.add(Container(color: Colors.white));
  }

  Widget createStones(int rowI) {
    var row = Row(
      mainAxisAlignment: MainAxisAlignment.end,
      children: <Widget>[
        Container(color: Colors.grey)
      ],
    );

    addParking(row, rowI);
    addStored(row, rowI);

    return row;
  }

  @override
  Widget build(BuildContext context) {
    return new Scaffold(
      body: new OrientationBuilder(
        builder: (context, orientation) {
          var background = Center(
            child: new Container(
              child: RTCVideoView(_localRenderer),
            ),
          );

          var parking = Column(
            children: <Widget>[
              Expanded(
                child: createStones(0),
              ),
              Expanded(
                child: createStones(1),
              ),
              Expanded(
                child: createStones(2),
              ),
              Expanded(
                child: createStones(3),
              ),
              Expanded(
                child: createStones(4),
              )
            ]
          );

          var foreground = Row(
            children: <Widget>[
              Expanded(
                child: Text('Left'),
              ),
              Expanded(
                child: parking,
              )
            ],
          );

          return new Material(
            child: Stack(
              children: <Widget>[
                background,
                foreground,
              ],
            ),
          );
        },
      ),
    );
  }
}