// Replace the code in main.dart with the following.

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'src/route_item.dart';
import 'src/start_screen.dart';
import 'package:flutter_webrtc/webrtc.dart';

void main(){
  runApp(new BlauApp());
}

class BlauApp extends StatefulWidget {

  @override
  _BlauAppState createState() => new _BlauAppState();
}


class _BlauAppState extends State<BlauApp> {
  List<RouteItem> items;

  @override
  initState() {
    super.initState();
    _initItems();
  }

  _initItems() {
    items = <RouteItem>[
    ];
  }

  @override
  Widget build(BuildContext context) {
    return new MaterialApp(
      title: "Select Players",
      home: new StartScreen(),
    );
  }
}

