import 'package:flutter/material.dart';
import 'package:go_todo/data/gt_api.dart';
import 'package:go_todo/presentation/route_scaffold.dart';
import 'package:go_todo/src/show_error_snack.dart';

import '../../src/get_snack_bar.dart';
import '../../widgets/gt_loading_button.dart';
import '../../widgets/gt_small_width_container.dart';
import '../../widgets/gt_text_field.dart';

class RegisterRoute extends StatefulWidget {
  const RegisterRoute({super.key});

  @override
  State<RegisterRoute> createState() => _RegisterRouteState();
}

class _RegisterRouteState extends State<RegisterRoute> {
  final usernameController = TextEditingController();
  final passwordController = TextEditingController();
  final passwordVerifyController = TextEditingController();
  final formKey = GlobalKey<FormState>();
  var isLoading = false;

  @override
  void dispose() {
    super.dispose();
    usernameController.dispose();
    passwordController.dispose();
    passwordVerifyController.dispose();
  }

  Future<void> register(BuildContext context) async {
    if (!formKey.currentState!.validate()) {
      return;
    }

    setState(() {
      isLoading = true;
    });
    final username = usernameController.text.trim();
    final password = passwordController.text;

    try {
      await GtApi().createUser(username, password);
      if (context.mounted) {
        final snackBar = getSnackBar(
          context: context,
          content: Text(
            'Successfully created account for $username',
          ),
        );
        ScaffoldMessenger.of(context).clearSnackBars();
        ScaffoldMessenger.of(context).showSnackBar(snackBar);
        Navigator.pop(context);
      }
    } on GtApiException catch (error) {
      if (context.mounted) {
        showErrorSnack(context, error, {
          GtApiExceptionType.conflict:
              'User with username $username already exists.',
        });
      }
    } catch (error) {
      if (context.mounted) {
        final snackBar = getSnackBar(
          context: context,
          content: const Text('There was error totally not handled...'),
        );
        ScaffoldMessenger.of(context).clearSnackBars();
        ScaffoldMessenger.of(context).showSnackBar(snackBar);
      }
    } finally {
      setState(() {
        isLoading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return RouteScaffold(
      showDrawer: false,
      body: GtSmallWidthContainer(
        child: Form(
          key: formKey,
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: <Widget>[
              Text(
                'Create Account',
                style: Theme.of(context).textTheme.headlineMedium,
              ),
              Padding(
                padding: const EdgeInsets.symmetric(vertical: 10),
                child: GtTextField(
                  controller: usernameController,
                  textInputAction: TextInputAction.next,
                  filled: true,
                  label: 'Username',
                  hint: 'Paroni, Julma-Hurtta, Liisa...',
                  leading: const Icon(Icons.person),
                  autovalidateMode: AutovalidateMode.onUserInteraction,
                  validator: (String? value) {
                    if (value == null || value.isEmpty) {
                      return "Username can't be empty.";
                    }
                    var re = RegExp(r'[^\p{L}\p{N}\s_-]', unicode: true);
                    if (re.hasMatch(value)) {
                      return 'Username contains disallowed characters.';
                    }
                    if (value.length < 3) {
                      return 'Username is too short.';
                    } else if (value.length > 20) {
                      return 'Username is too long.';
                    }
                    return null;
                  },
                ),
              ),
              Padding(
                padding: const EdgeInsets.symmetric(vertical: 10.0),
                child: GtTextField(
                  controller: passwordController,
                  textInputAction: TextInputAction.next,
                  filled: true,
                  label: 'Password',
                  isSecret: true,
                  leading: const Icon(Icons.password),
                  autovalidateMode: AutovalidateMode.onUserInteraction,
                  validator: (String? value) {
                    if (value == null || value.isEmpty) {
                      return "Password can't be empty";
                    }
                    final reLetter = RegExp(r'\p{L}', unicode: true);
                    final reNumber = RegExp(r'\d');
                    final reSpecialChar = RegExp(r'\p{P}|\p{S}', unicode: true);
                    if (!reLetter.hasMatch(value)) {
                      return 'Has to contain a letter';
                    }
                    if (!reNumber.hasMatch(value)) {
                      return 'Has to contain a number';
                    }
                    if (!reSpecialChar.hasMatch(value)) {
                      return 'Has to contain a special character';
                    }
                    if (value.length < 8) {
                      return 'Too short';
                    } else if (value.length > 32) {
                      return 'Too long';
                    }
                    return null;
                  },
                ),
              ),
              Padding(
                padding: const EdgeInsets.symmetric(vertical: 10.0),
                child: GtTextField(
                  controller: passwordVerifyController,
                  textInputAction: TextInputAction.done,
                  filled: true,
                  label: 'Confirm Password',
                  isSecret: true,
                  leading: const Icon(
                    Icons.password,
                  ),
                  autovalidateMode: AutovalidateMode.onUserInteraction,
                  validator: (String? value) {
                    if (value == null || value.isEmpty) {
                      return 'Does not match';
                    }
                    if (value != passwordController.value.text) {
                      return 'Does not match';
                    }
                    return null;
                  },
                ),
              ),
              Padding(
                padding: const EdgeInsets.only(top: 10.0),
                child: SizedBox(
                  width: double.infinity,
                  child: GtLoadingButton(
                    isLoading: isLoading,
                    onPressed: () async => await register(context),
                    text: 'Create',
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
