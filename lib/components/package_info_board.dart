// Copyright 2023 FMNX Linux team.
// This code is covered by GPL license, which can be found in LICENSE file.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io
import 'package:repo/components/fmnx_button.dart';
import 'package:repo/components/notification.dart';
import 'package:repo/constants.dart';
import 'package:repo/generated/v1/pack.pbgrpc.dart';
import 'package:data_table_2/data_table_2.dart';
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
            child: DataTable2(
              columnSpacing: defaultPadding,
              minWidth: 600,
              columns: [
                DataColumn2(
                  label: Text("Property"),
                  size: ColumnSize.S,
                ),
                DataColumn2(
                  label: Text("Value"),
                  size: ColumnSize.L,
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
                  key: "Architecture",
                  value: widget.description.architecture,
                  icon: Icons.computer,
                  color: Colors.orange,
                ),
                formDataRow(
                  value: widget.description.installedSize,
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
                  key: "Licenses",
                  value: widget.description.licenses,
                  icon: Icons.assignment,
                  color: Colors.purple,
                ),
                formDataRow(
                  value: widget.description.groups,
                  key: "Groups",
                  icon: Icons.apps,
                  color: Colors.brown,
                ),
                formDataRow(
                  value: widget.description.provides,
                  key: "Provides",
                  icon: Icons.crop_free,
                  color: Colors.indigo,
                ),
                formDataRow(
                  value: widget.description.packager,
                  key: "Packager",
                  icon: Icons.person,
                  color: Colors.blue,
                ),
                formDataRow(
                  value: widget.description.buildDate,
                  key: "Build date",
                  icon: Icons.donut_small,
                  color: Colors.cyan,
                ),
                formDataRow(
                  value: widget.description.conflictsWith,
                  key: "Conflicts",
                  icon: Icons.do_disturb,
                  color: Colors.red,
                ),
                formDataRow(
                  value: widget.description.replaces,
                  key: "Replaces",
                  icon: Icons.crop_3_2,
                  color: Colors.red,
                ),
                formDataRow(
                  value: widget.description.installScript,
                  key: "Installation script",
                  icon: Icons.code,
                  color: Colors.deepPurple,
                ),
                formDataRow(
                  value: widget.description.validatedBy,
                  key: "Validated by",
                  icon: Icons.check_circle,
                  color: Colors.green,
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
        Row(
          mainAxisAlignment: MainAxisAlignment.start,
          children: [
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
            Padding(padding: EdgeInsets.all(defaultPadding)),
            Text(key),
          ],
        ),
      ),
      DataCell(Text(value), onTap: () {
        if (key == "Official web URL") {
          launchUrl(Uri.parse(value));
        }
      }),
    ],
  );
}
