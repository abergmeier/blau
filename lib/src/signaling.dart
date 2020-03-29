
import 'package:flutter_webrtc/webrtc.dart';
import 'websocket.dart';

class Signaling {

  SimpleWebSocket _socket;
  var _host;
  var _port = 8086;

  var _peerConnections = new Map<String, RTCPeerConnection>();
  var _dataChannels = new Map<String, RTCDataChannel>();
  var _turnCredential;

  Map<String, dynamic> _iceServers = {
    'iceServers': [
      {'url': 'stun:stun.l.google.com:19302'},
      /*
       * turn server configuration example.
      {
        'url': 'turn:123.45.67.89:3478',
        'username': 'change_to_real_user',
        'credential': 'change_to_real_secret'
      },
       */
    ]
  };

  Signaling(this._host);

  void connect() async {
    var url = 'https://$_host:$_port/ws';
    _socket = SimpleWebSocket(url);

    print('connect to $url');

    _socket.onOpen = () {
      print('onOpen');
    };

    _socket.onMessage = (message) {
      print('Recivied data: ' + message);
    };

    _socket.onClose = (int code, String reason) {
      print('Closed by server [$code => $reason]!');
   };

    await _socket.connect();
  }
}