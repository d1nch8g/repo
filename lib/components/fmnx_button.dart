import 'package:fmnxpkg/generated/v1/pacman.pb.dart';
import 'package:fmnxpkg/responsive.dart';
import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../../constants.dart';

class FmnxButton extends StatefulWidget {
  final String text;
  final void Function()? onPressed;
  final IconData icon;
  const FmnxButton({
    Key? key,
    this.onPressed,
    required this.text,
    required this.icon,
  }) : super(key: key);

  @override
  State<FmnxButton> createState() => _FmnxButtonState();
}

class _FmnxButtonState extends State<FmnxButton> {
  Widget buttonWidget = Container();

  checkAuth() async {
    var prefs = await SharedPreferences.getInstance();
    var checked = await stub.checkToken(CheckTokenRequest(
      token: prefs.getString("token"),
    ));
    if (checked.upToDate ||
        widget.text == "Authorize" ||
        widget.text == "Login" ||
        widget.text == "Close") {
      setState(() {
        buttonWidget = ElevatedButton.icon(
          style: TextButton.styleFrom(
            padding: EdgeInsets.symmetric(
              horizontal:
                  defaultPadding * (Responsive.isMobile(context) ? 0.6 : 1.5),
              vertical:
                  defaultPadding / (Responsive.isMobile(context) ? 1.25 : 1),
            ),
            backgroundColor: primaryColor,
          ),
          onPressed: widget.onPressed,
          icon: Icon(widget.icon),
          label: Text(widget.text),
        );
      });
    }
  }

  @override
  void initState() {
    super.initState();
    checkAuth();
  }

  @override
  Widget build(BuildContext context) {
    return buttonWidget;
  }
}
