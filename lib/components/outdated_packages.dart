// Copyright 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

import 'package:repo/generated/v1/pack.pb.dart';
import 'package:repo/components/fmnx_button.dart';
import 'package:repo/components/update_packages.dart';
import 'package:flutter/material.dart';

import '../../../constants.dart';
import 'package_chart.dart';
import 'outdated_package_card.dart';

class OutdatedPackages extends StatelessWidget {
  final double outdated;
  final double total;
  final List<OutdatedPackage> outdatedPackagesList;
  final void Function() updateCallback;
  const OutdatedPackages({
    Key? key,
    required this.outdated,
    required this.total,
    required this.outdatedPackagesList,
    required this.updateCallback,
  }) : super(key: key);

  List<OutdatedPackageCard> convertPackages(
      List<OutdatedPackage> outdatedPackagesList) {
    List<OutdatedPackageCard> output = [];
    outdatedPackagesList.forEach((element) {
      output.add(OutdatedPackageCard(
        name: element.name,
        latestVersion: element.latestVersion,
        currVersion: element.currentVersion,
      ));
    });
    return output;
  }

  double getHeight(double outdated) {
    if (outdated == 1.0) {
      return 94;
    }
    if (outdated == 2.0) {
      return 188;
    }
    if (outdated == 3.0) {
      return 282;
    }
    return 396;
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
          Text(
            "Outdated packages",
            style: TextStyle(
              fontSize: 18,
              fontWeight: FontWeight.w500,
            ),
          ),
          SizedBox(height: defaultPadding),
          PackageChart(
            outdated: outdated,
            uptodate: total - outdated,
          ),
          if (outdated == 0)
            Container(
              height: getHeight(1),
              child: Center(
                child: Text(
                  "Evrything up to date",
                  style: TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.w400,
                  ),
                ),
              ),
            ),
          if (outdated > 0)
            Container(
              height: getHeight(outdated),
              child: ListView(
                children: convertPackages(outdatedPackagesList),
              ),
            ),
          if (outdated > 0)
            Center(
              child: Padding(
                padding: const EdgeInsets.all(defaultPadding),
                child: FmnxButton(
                  text: "Update",
                  icon: Icons.refresh,
                  onPressed: () {
                    showUpdateNotification(context, updateCallback);
                  },
                ),
              ),
            ),
        ],
      ),
    );
  }
}
