// Copyright 2023 FMNX Linux team.
// This code is covered by GPL license, which can be found in LICENSE file.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io
import 'package:fl_chart/fl_chart.dart';
import 'package:flutter/material.dart';

import '../../../constants.dart';

class PackageChart extends StatelessWidget {
  final double outdated;
  final double uptodate;
  const PackageChart({
    Key? key,
    required this.outdated,
    required this.uptodate,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: 200,
      child: Stack(
        children: [
          PieChart(
            PieChartData(
              sectionsSpace: 0,
              centerSpaceRadius: 70,
              startDegreeOffset: -90,
              sections: [
                PieChartSectionData(
                  color: Color.fromARGB(255, 64, 138, 250),
                  value: uptodate,
                  showTitle: false,
                  radius: 22,
                ),
                PieChartSectionData(
                  color: Color(0xFFEE2727),
                  value: outdated,
                  showTitle: false,
                  radius: 25,
                ),
              ],
            ),
          ),
          Positioned.fill(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                SizedBox(height: defaultPadding),
                Text(
                  "$outdated",
                  style: Theme.of(context).textTheme.headlineMedium!.copyWith(
                        color: Colors.white,
                        fontWeight: FontWeight.w600,
                        height: 0.5,
                      ),
                ),
                Text("of $uptodate")
              ],
            ),
          ),
        ],
      ),
    );
  }
}
