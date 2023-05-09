// Copyright 2023 FMNX Linux team.
// This code is covered by GPL license, which can be found in LICENSE file.
// Additional information could be found on official web page: https://fmnx.io/
// Email: help@fmnx.io
import 'package:flutter/material.dart';

import '../../../constants.dart';

class NotificationPopup extends StatefulWidget {
  final String message;
  final IconData icon;
  final Duration duration;
  const NotificationPopup({
    required this.message,
    required this.icon,
    required this.duration,
  });

  @override
  State<NotificationPopup> createState() => _NotificationPopupState();
}

class _NotificationPopupState extends State<NotificationPopup> {
  @override
  void initState() {
    super.initState();
    Future.delayed(widget.duration, () {
      Navigator.pop(context);
    });
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      height: 54,
      decoration: BoxDecoration(
        color: secondaryColor,
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.35),
            blurRadius: 1,
            offset: const Offset(-1, -1),
          ),
        ],
      ),
      child: Padding(
        padding: const EdgeInsets.all(8),
        child: AnimatedSwitcher(
          duration: const Duration(milliseconds: 458),
          child: Row(
            key: UniqueKey(),
            mainAxisSize: MainAxisSize.max,
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(
                widget.icon,
                size: 24,
                color: Colors.white,
              ),
              const SizedBox(width: 8),
              Text(
                widget.message,
                style: const TextStyle(
                  fontSize: 16,
                  color: Colors.white,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
