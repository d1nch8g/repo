import 'package:file_picker/file_picker.dart';
import 'package:flutter/material.dart';
import 'package:fmnxpkg/components/notification.dart';
import 'package:fmnxpkg/constants.dart';
import 'package:fmnxpkg/generated/v1/pacman.pb.dart';
import 'package:shared_preferences/shared_preferences.dart';

uploadFile(BuildContext context) async {
  try {
    var prefs = await SharedPreferences.getInstance();
    var pick = await FilePicker.platform.pickFiles(
      type: FileType.custom,
      allowMultiple: false,
      onFileLoading: (status) => print(status),
      allowedExtensions: ['pkg.tar.zst'],
    );
    if (pick == null) {
      return;
    }
    showBottomSheet(
      context: context,
      builder: (BuildContext context) {
        return NotificationPopup(
          message: "Upload have started...",
          icon: Icons.warning_amber_rounded,
          duration: Duration(milliseconds: 2342),
        );
      },
    );
    await stub.upload(UploadRequest(
      content: pick.files.first.bytes,
      name: pick.files.first.name,
      token: prefs.getString("token"),
    ));
    showBottomSheet(
      context: context,
      builder: (BuildContext context) {
        return NotificationPopup(
          message: "Done!",
          icon: Icons.check_circle_outline,
          duration: Duration(milliseconds: 1342),
        );
      },
    );
  } catch (e) {
    print(e);
    showBottomSheet(
      context: context,
      builder: (BuildContext context) {
        return NotificationPopup(
          message: "Upload failed",
          icon: Icons.error,
          duration: Duration(milliseconds: 1342),
        );
      },
    );
  }
}
