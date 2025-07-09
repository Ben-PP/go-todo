import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_todo/application/authentication_provider.dart';
import 'package:go_todo/data/gt_api.dart';
import 'package:go_todo/presentation/route_scaffold.dart';
import 'package:go_todo/src/get_snack_bar.dart';
import 'package:go_todo/widgets/gt_loading_button.dart';
import 'package:go_todo/widgets/gt_small_width_container.dart';
import 'package:go_todo/widgets/gt_text_field.dart';

class ApiUrlRoute extends ConsumerStatefulWidget {
  const ApiUrlRoute({super.key, this.canCancel = false});
  final bool canCancel;

  @override
  ConsumerState<ApiUrlRoute> createState() => _ApiUrlRouteState();
}

class _ApiUrlRouteState extends ConsumerState<ApiUrlRoute> {
  final TextEditingController _apiUrlController = TextEditingController();
  var isLoading = false;
  var isSaveDisabled = true;
  final formKey = GlobalKey<FormState>();

  @override
  void dispose() {
    super.dispose();
    _apiUrlController.dispose();
  }

  String? validateUrl(String? value) {
    if (value == null) {
      return "Url can't be empty";
    }
    if (!value.startsWith('http://') && !value.startsWith('https://')) {
      return 'Url must start with http:// or https://';
    }
    return null;
  }

  var snackMessage = '';
  var snackIsError = false;

  Future<void> saveUrl(BuildContext context) async {
    if (!formKey.currentState!.validate()) {
      return;
    }
    setState(() {
      isLoading = true;
    });
    try {
      final url = _apiUrlController.text.trim();
      await GtApi().setBaseUrl(url);
      await ref.read(authenticationProvider.notifier).refresh();
      if (context.mounted) {
        snackIsError = false;
        snackMessage = 'Added $url as Base URL!';
        Navigator.pop(context);
        return;
      }
      snackIsError = true;
      snackMessage = 'Context was not mounted. Weird...';
    } on GtApiException catch (gtError) {
      snackIsError = true;
      switch (gtError.type) {
        case GtApiExceptionType.hostUnknown:
          snackMessage = 'Host did not appear to exist.';
          break;
        case GtApiExceptionType.hostNotResponding:
          snackMessage = 'Host refused connection. Feels kinda sad.';
          break;
        case GtApiExceptionType.malformedBody:
          snackMessage = 'Refresh requests body was malformed.';
        case GtApiExceptionType.refreshJWTUnauthorized:
          snackMessage = gtError.cause;
        case GtApiExceptionType.unknown:
          snackMessage = 'Confusion of the highest order! Unknown error :(';
        case GtApiExceptionType.urlNull:
          snackMessage = 'Could not set url.';
        case GtApiExceptionType.serverError:
          snackMessage = 'Server just rage quit on you. 500 all the way!';
        default:
          snackMessage = 'This should never be shown';
      }
    } catch (error) {
      snackMessage = 'Encountered unhandled error :(';
      snackIsError = true;
      return;
    } finally {
      setState(() {
        isLoading = false;
      });
      if (context.mounted) {
        var snackBar = getSnackBar(
          context: context,
          content: Text(snackMessage),
          isError: snackIsError,
        );
        ScaffoldMessenger.of(context).showSnackBar(snackBar);
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return PopScope(
      canPop: widget.canCancel,
      child: RouteScaffold(
        showDrawer: false,
        implyLeading: widget.canCancel,
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
              Form(
                key: formKey,
                child: GtTextField(
                  label: 'API Host URL',
                  hint: 'https://api.example.com:8000',
                  controller: _apiUrlController,
                  filled: true,
                  autovalidateMode: AutovalidateMode.disabled,
                  validator: validateUrl,
                ),
              ),
              const SizedBox(height: 20),
              Row(
                children: [
                  if (widget.canCancel)
                    Expanded(
                      child: GtLoadingButton(
                          text: 'Cancel',
                          onPressed: () => Navigator.pop(context)),
                    ),
                  if (widget.canCancel)
                    const SizedBox(
                      width: 15,
                    ),
                  Expanded(
                    child: GtLoadingButton(
                      onPressed: () => saveUrl(context),
                      text: 'Save',
                      isLoading: isLoading,
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}
