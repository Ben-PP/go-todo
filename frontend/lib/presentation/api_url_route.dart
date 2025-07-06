import 'package:flutter/material.dart';
import 'package:go_todo/data/gt_api.dart';
import 'package:go_todo/presentation/route_scaffold.dart';
import 'package:go_todo/widgets/gt_small_width_container.dart';
import 'package:go_todo/widgets/gt_text_field.dart';

class ApiUrlRoute extends StatefulWidget {
  const ApiUrlRoute({super.key, this.onSuccess});
  final Function? onSuccess;

  @override
  State<ApiUrlRoute> createState() => _ApiUrlRouteState();
}

class _ApiUrlRouteState extends State<ApiUrlRoute> {
  final TextEditingController _apiUrlController = TextEditingController();

  @override
  void dispose() {
    super.dispose();
    _apiUrlController.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return RouteScaffold(
      title: const Text('Set up API Host'),
      body: GtSmallWidthContainer(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(
              'Enter Address',
              style: Theme.of(context).textTheme.headlineMedium,
            ),
            Text(
              'This is the address of your Go Todo API server.',
              style: Theme.of(context).textTheme.bodyMedium,
            ),
            GtTextField(
              label: 'API Host URL',
              hint: 'https://api.example.com:8000',
              controller: _apiUrlController,
              filled: true,
            ),
            const SizedBox(height: 20),
            SizedBox(
              width: double.infinity,
              child: ElevatedButton(
                  onPressed: () async {
                    if (_apiUrlController.text.trim().isEmpty) {
                      return;
                    }
                    GtApi()
                        .setBaseUrl(_apiUrlController.text.trim())
                        .then((value) {
                      if (widget.onSuccess != null) widget.onSuccess!();
                    }).onError((error, _) {});
                  },
                  child: const Text('Save')),
            ),
          ],
        ),
      ),
    );
  }
}
