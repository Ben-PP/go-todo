import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_todo/application/authentication_provider.dart';
import 'package:go_todo/data/gt_api.dart';
import 'package:go_todo/presentation/auth_view/register_route.dart';
import 'package:go_todo/src/create_gt_route.dart';
import 'package:go_todo/src/get_snack_bar.dart';
import 'package:go_todo/src/show_error_snack.dart';
import 'package:go_todo/widgets/gt_loading_button.dart';
import 'package:go_todo/widgets/gt_small_width_container.dart';
import 'package:go_todo/widgets/gt_text_field.dart';

class LoginView extends ConsumerStatefulWidget {
  const LoginView({super.key});

  @override
  ConsumerState<LoginView> createState() => _LoginViewState();
}

class _LoginViewState extends ConsumerState<LoginView> {
  final usernameController = TextEditingController();
  final passwordController = TextEditingController();

  var isLoading = false;

  Future<void> login(BuildContext context) async {
    var uname = usernameController.text.trim();
    var passwd = passwordController.text.trim();
    if (passwd.isEmpty || uname.isEmpty) {
      final snackBar = getSnackBar(
        context: context,
        content: Text(
          'Empty ${uname.isEmpty ? 'username' : 'password'}!',
        ),
        isError: true,
      );
      ScaffoldMessenger.of(context).clearSnackBars();
      ScaffoldMessenger.of(context).showSnackBar(snackBar);
      return;
    }
    setState(() {
      isLoading = true;
    });
    try {
      await ref.read(authenticationProvider.notifier).login(uname, passwd);
      if (context.mounted) {
        ScaffoldMessenger.of(context).clearSnackBars();
        ScaffoldMessenger.of(context).showSnackBar(
          getSnackBar(
            context: context,
            content: Text('Successfully logged in as $uname'),
          ),
        );
      }
    } on GtApiException catch (error) {
      if (context.mounted) {
        showErrorSnack(context, error, map: {
          GtApiExceptionType.malformedBody:
              'Login requests body was malformed.',
          GtApiExceptionType.unauthorized: "Username/Password doesn't match.",
          GtApiExceptionType.serverError:
              'You broke the server (500) :(\nContact your personal support guy.',
          GtApiExceptionType.unknownResponse:
              'Something mysterious was handled incorrectly...',
          GtApiExceptionType.hostNotResponding:
              'Your server is not talking to us.',
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
      setState(() {
        isLoading = false;
      });
    }
  }

  @override
  void dispose() {
    super.dispose();
    usernameController.dispose();
    passwordController.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return GtSmallWidthContainer(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: <Widget>[
          Text(
            'Login',
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
            ),
          ),
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 10.0),
            child: GtTextField(
              controller: passwordController,
              textInputAction: TextInputAction.done,
              filled: true,
              label: 'Password',
              isSecret: true,
              leading: const Icon(Icons.password),
            ),
          ),
          Padding(
            padding: const EdgeInsets.only(top: 10.0),
            child: SizedBox(
              width: double.infinity,
              child: GtLoadingButton(
                isLoading: isLoading,
                onPressed: () async => await login(context),
                text: 'Login',
              ),
            ),
          ),
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 10),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                const Text("Don't have an account?"),
                const SizedBox(
                  width: 2,
                ),
                TextButton(
                    onPressed: () {
                      Navigator.push(
                        context,
                        createGtRoute(context, const RegisterRoute()),
                      );
                    },
                    child: const Text('Register')),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
