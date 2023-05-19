// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

import 'package:repo/components/fmnx_button.dart';
import 'package:repo/components/notification.dart';
import 'package:repo/constants.dart';
import 'package:repo/generated//pack.pbgrpc.dart';
import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:url_launcher/url_launcher.dart';

class PackageInfoBoard extends StatefulWidget {
  final DescribeResponse description;

  const PackageInfoBoard({
    Key? key,
    required this.description,
  }) : super(key: key);

  @override
  State<PackageInfoBoard> createState() => _PackageInfoBoardState();
}

class _PackageInfoBoardState extends State<PackageInfoBoard> {
  Widget removeWidget = Container();

  setRemoveWidget() async {
    var prefs = await SharedPreferences.getInstance();
    var checked = await stub.checkToken(CheckTokenRequest(
      token: prefs.getString("token") ?? "",
    ));
    if (checked.upToDate) {
      setState(() {
        removeWidget = Center(
          child: Padding(
            padding: const EdgeInsets.all(8.0),
            child: FmnxButton(
              text: "delete",
              icon: Icons.do_disturb,
              onPressed: () async {
                var prefs = await SharedPreferences.getInstance();
                try {
                  await stub.remove(RemoveRequest(
                    package: widget.description.name,
                    token: prefs.getString("token"),
                  ));
                  showBottomSheet(
                    context: context,
                    builder: (BuildContext context) {
                      return NotificationPopup(
                        message: "Package removed",
                        icon: Icons.check_circle_outline_rounded,
                        duration: Duration(milliseconds: 2342),
                      );
                    },
                  );
                } catch (e) {
                  showBottomSheet(
                    context: context,
                    builder: (BuildContext context) {
                      return NotificationPopup(
                        message: e.toString(),
                        icon: Icons.error,
                        duration: Duration(milliseconds: 2342),
                      );
                    },
                  );
                }
              },
            ),
          ),
        );
      });
    }
  }

  @override
  void initState() {
    super.initState();
    setRemoveWidget();
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: EdgeInsets.all(defaultPadding),
      decoration: BoxDecoration(
        color: secondaryColor,
        borderRadius: const BorderRadius.all(Radius.circular(10)),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Center(
            child: Text(
              widget.description.name.toUpperCase(),
              style: Theme.of(context).textTheme.titleMedium,
            ),
          ),
          SizedBox(
            width: double.infinity,
            child: DataTable(
              columnSpacing: defaultPadding,
              columns: [
                DataColumn(
                  label: Text(""),
                ),
                DataColumn(
                  label: Text("Property"),
                ),
                DataColumn(
                  label: Text("Value"),
                ),
              ],
              rows: [
                formDataRow(
                  key: "Current version",
                  value: widget.description.version,
                  icon: Icons.receipt,
                  color: Colors.blueGrey,
                ),
                formDataRow(
                  key: "Description",
                  value: widget.description.description,
                  icon: Icons.article,
                  color: Colors.deepOrange,
                ),
                formDataRow(
                  value: widget.description.size,
                  key: "Storage size",
                  icon: Icons.sd_card,
                  color: Colors.teal,
                ),
                formDataRow(
                  key: "Official web URL",
                  value: widget.description.url,
                  icon: Icons.web,
                  color: Colors.grey,
                ),
                formDataRow(
                  value: widget.description.buildDate,
                  key: "Build date",
                  icon: Icons.donut_small,
                  color: Colors.cyan,
                ),
                formDataRow(
                  value: widget.description.packName,
                  key: "Pack link",
                  icon: Icons.broadcast_on_home,
                  color: Colors.orange,
                ),
                formDataRow(
                  value: widget.description.packVersion,
                  key: "Pack version",
                  icon: Icons.browser_updated,
                  color: Colors.purple,
                ),
                formDataRow(
                  value: widget.description.packBranch,
                  key: "Default branch",
                  icon: Icons.arrow_back_ios_new_outlined,
                  color: Colors.blue,
                ),
                formDataRow(
                  value: widget.description.dependsOn,
                  key: "Depends on",
                  icon: Icons.lock,
                  color: Colors.amber,
                ),
                formDataRow(
                  value: widget.description.requiredBy,
                  key: "Required by",
                  icon: Icons.lock_outline,
                  color: Colors.red,
                ),
              ],
            ),
          ),
          removeWidget,
        ],
      ),
    );
  }
}

DataRow formDataRow({
  required String key,
  required String value,
  required IconData icon,
  required Color color,
}) {
  return DataRow(
    cells: [
      DataCell(
        Container(
          decoration: BoxDecoration(
            color: color,
            borderRadius: const BorderRadius.all(Radius.circular(4)),
          ),
          child: Padding(
            padding: const EdgeInsets.all(2.0),
            child: Icon(
              icon,
              size: 18,
              color: Colors.white,
            ),
          ),
        ),
      ),
      DataCell(
        Padding(
          padding: const EdgeInsets.all(2.0),
          child: SizedBox(
            width: 120,
            child: Text(key),
          ),
        ),
      ),
      DataCell(
        Text(value + " " * 1000),
        onTap: () {
          if (key == "Official web URL") {
            launchUrl(Uri.parse(value));
          }
        },
      ),
    ],
  );
}
