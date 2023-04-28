import 'package:repo/generated/v1/pacman.pbgrpc.dart';
import 'package:flutter/material.dart';
import 'package:grpc/grpc_web.dart';

const primaryColor = Color.fromARGB(255, 116, 116, 116);
const secondaryColor = Color.fromARGB(255, 54, 54, 54);
const backgroundColor = Color.fromARGB(255, 36, 36, 36);

const defaultPadding = 16.0;

var chan = GrpcWebClientChannel.xhr(Uri.parse("http://localhost:8080/"));
var stub = PacmanServiceClient(chan, CallOptions(timeout: Duration(days: 1)));
