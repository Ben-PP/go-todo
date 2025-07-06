import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_todo/application/authentication_provider.dart';
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

  // TODO Add registration view
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
              child: ElevatedButton(
                  onPressed: () async {
                    // TODO Add validation
                    var uname = usernameController.text.trim();
                    var passwd = passwordController.text.trim();
                    await ref
                        .read(authenticationProvider.notifier)
                        .login(uname, passwd);
                  },
                  child: const Text('Login')),
            ),
          ),
        ],
      ),
    );
  }
}
