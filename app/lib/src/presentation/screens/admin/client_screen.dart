import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:thetiptop_client/src/presentation/widgets/layouts/layout_admin_widget.dart';

class ClientScreen extends HookConsumerWidget {
  const ClientScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final formKey = GlobalKey<FormState>();

    return LayoutAdminWidget(
      child: Text("ok"),
    );
  }
}
