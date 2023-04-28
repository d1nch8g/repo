import 'package:file_picker/file_picker.dart';
import 'package:flutter/material.dart';
import 'package:flutter_spinkit/flutter_spinkit.dart';
import 'package:repo/constants.dart';
import 'package:repo/generated/v1/pacman.pb.dart';
import 'package:shared_preferences/shared_preferences.dart';

uploadFile(BuildContext context) async {
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
      return UploadContent(
        uploadCallback: () async {
          try {
            var prefs = await SharedPreferences.getInstance();
            await stub.upload(UploadRequest(
              content: pick.files.first.bytes,
              name: pick.files.first.name,
              token: prefs.getString("token"),
            ));
            return true;
          } catch (e) {
            return false;
          }
        },
      );
    },
  );
}

class UploadContent extends StatefulWidget {
  final Future<bool> Function() uploadCallback;
  const UploadContent({Key? key, required this.uploadCallback})
      : super(key: key);

  @override
  State<UploadContent> createState() => _UploadContentState();
}

class _UploadContentState extends State<UploadContent> {
  Widget currIcon = SpinKitDoubleBounce(
    size: 24,
    color: Colors.white,
  );
  String currMessage = "Uploading and installing...";

  upload() async {
    var success = await widget.uploadCallback();
    if (success) {
      setState(() {
        currIcon = Icon(
          Icons.check_circle,
          size: 24,
          color: Colors.white,
        );
        currMessage = "Upload success.";
        Future.delayed(Duration(milliseconds: 1327), () {
          Navigator.pop(context);
        });
      });
    } else {
      setState(() {
        currIcon = Icon(
          Icons.error,
          size: 24,
          color: Colors.white,
        );
        currMessage = "Error occured.";
        Future.delayed(Duration(milliseconds: 1327), () {
          Navigator.pop(context);
        });
      });
    }
  }

  @override
  void initState() {
    super.initState();
    upload();
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
              currIcon,
              const SizedBox(width: 8),
              Text(
                currMessage,
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
