import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_todo/application/todo_list.dart';
import 'package:go_todo/data/gt_api.dart';
import 'package:go_todo/presentation/route_scaffold.dart';
import 'package:go_todo/src/get_snack_bar.dart';
import 'package:go_todo/src/show_error_snack.dart';
import 'package:go_todo/widgets/gt_loading_button.dart';
import 'package:go_todo/widgets/gt_small_width_container.dart';
import 'package:go_todo/widgets/gt_text_field.dart';

class CreateListRoute extends ConsumerStatefulWidget {
  const CreateListRoute({super.key});

  @override
  ConsumerState<CreateListRoute> createState() => _CreateListRouteState();
}

class _CreateListRouteState extends ConsumerState<CreateListRoute> {
  final titleController = TextEditingController();
  final descriptionController = TextEditingController();
  var isLoading = false;

  @override
  void dispose() {
    super.dispose();
    titleController.dispose();
    descriptionController.dispose();
  }

  createList(BuildContext context) async {
    final title = titleController.text.trim();
    final description = descriptionController.text.trim();
    if (title.isEmpty) {
      final snackBar = getSnackBar(
          context: context, content: const Text('Title cannot be empty'));
      ScaffoldMessenger.of(context).showSnackBar(snackBar);
      return;
    }

    setState(() {
      isLoading = true;
    });
    try {
      await ref.read(todoListProvider.notifier).createList(
          title: title, description: description.isEmpty ? null : description);
      if (context.mounted) {
        ScaffoldMessenger.of(context).clearSnackBars();
        ScaffoldMessenger.of(context).showSnackBar(
          getSnackBar(
            context: context,
            content: const Text('New list created!'),
          ),
        );
        Navigator.of(context).pop();
      }
    } on GtApiException catch (error) {
      if (context.mounted) {
        showErrorSnack(context, error, map: {
          GtApiExceptionType.malformedBody:
              'Something was odd in the request. They say it was malformed.',
          GtApiExceptionType.unauthorized:
              'Hey! This action is not for unauthorized users.',
          GtApiExceptionType.serverError:
              'The server is having a bad day. Please try again later.',
          GtApiExceptionType.unknownResponse:
              'We do not know what the response meant.',
          GtApiExceptionType.hostNotResponding:
              'The server is not responding. Please check your connection.',
        });
      }
    } catch (error) {
      if (context.mounted) {
        ScaffoldMessenger.of(context).clearSnackBars();
        ScaffoldMessenger.of(context).showSnackBar(
          getSnackBar(
            context: context,
            content: const Text(
                'This error was not handled at all. Fix the thrash...'),
            isError: true,
          ),
        );
      }
    } finally {
      if (context.mounted) {
        setState(() {
          isLoading = false;
        });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return RouteScaffold(
      implyLeading: true,
      showDrawer: false,
      title: const Text('New TODO List'),
      body: GtSmallWidthContainer(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Padding(
              padding: const EdgeInsets.symmetric(vertical: 10.0),
              child: Text(
                'Details of the List',
                style: Theme.of(context).textTheme.headlineSmall,
              ),
            ),
            Form(
              child: Column(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(vertical: 10.0),
                    child: GtTextField(
                      controller: titleController,
                      maxLength: 40,
                      filled: true,
                      label: 'Title',
                      hint: 'Title goes here...',
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(vertical: 10.0),
                    child: GtTextField(
                      controller: descriptionController,
                      maxLength: 150,
                      filled: true,
                      label: 'Description',
                      hint: 'Description goes here...',
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(vertical: 10.0),
                    child: SizedBox(
                      width: double.infinity,
                      child: GtLoadingButton(
                        text: 'Create List',
                        onPressed: () => createList(context),
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
