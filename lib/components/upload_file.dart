import 'package:file_picker/file_picker.dart';
import 'package:flutter/material.dart';
import 'package:flutter_spinkit/flutter_spinkit.dart';
import 'package:fmnx_pkg/constants.dart';
import 'package:fmnx_pkg/generated/v1/pacman.pb.dart';
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
      return UploadContent(file: pick);
    },
  );
}

class UploadContent extends StatefulWidget {
  final FilePickerResult file;

  const UploadContent({Key? key, required this.file}) : super(key: key);
  @override
  State<UploadContent> createState() => _UploadContentState();
}

class _UploadContentState extends State<UploadContent> {
  Widget currIcon = SpinKitDoubleBounce(
    size: 24,
    color: Colors.white,
  );
  String currMessage = "Uploading and installing...";

  uploadUpdate() async {
    var prefs = await SharedPreferences.getInstance();
    try {
      await stub.upload(UploadRequest(
        content: widget.file.files.first.bytes,
        name: widget.file.files.first.name,
        token: prefs.getString("token"),
      ));
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
    } catch (e) {
      setState(() {
        currIcon = Icon(
          Icons.error,
          size: 24,
          color: Colors.white,
        );
        currMessage = e.toString();
        Future.delayed(Duration(milliseconds: 1327), () {
          Navigator.pop(context);
        });
      });
    }
  }

  @override
  void initState() {
    super.initState();
    uploadFile(context);
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
