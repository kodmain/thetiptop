import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:thetiptop_client/src/presentation/screens/admin/admin_signin_screen.dart';
import 'package:thetiptop_client/src/presentation/screens/admin/client_screen.dart';
import 'package:thetiptop_client/src/presentation/screens/client/game_screen.dart';
import 'package:thetiptop_client/src/presentation/screens/client/history_screen.dart';
import 'package:thetiptop_client/src/presentation/screens/client/password_renew_screen.dart';
import 'package:thetiptop_client/src/presentation/screens/client/password_reset_screen.dart';
import 'package:thetiptop_client/src/presentation/screens/client/profil_screen.dart';
import 'package:thetiptop_client/src/presentation/screens/client/signin_screen.dart';
import 'package:thetiptop_client/src/presentation/screens/client/signup_screen.dart';

class AppRouter {
  static const String clientRouteName = '/admin/client';
  static const String adminSigninRouteName = '/admin/signin';

  static const String profilRouteName = '/profil';
  static const String historyRouteName = '/history';
  static const String gameRouteName = '/game';
  static const String signinRouteName = '/signin';
  static const String signupRouteName = '/signup';
  static const String passwordRenewRouteName = '/password/renew';
  static const String passwordResetRouteName = '/password/reset';

  /// https://pub.dev/packages/go_router
  /// https://pub.dev/documentation/go_router/latest/index.html
  ///
  /// Renvoi la définition des routes de l'application
  ///
  /// Warning, Builder vs PageBuilder : PageBuilder nous permet d'avoir une
  /// animation de transition personnalisée via les customTransitionPage
  ///
  static GoRouter routes() {
    return GoRouter(
      debugLogDiagnostics: true,
      initialLocation: '/admin/client',
      routes: [
        GoRoute(
          name: signinRouteName,
          path: signinRouteName,
          pageBuilder: (BuildContext context, GoRouterState state) {
            return CustomTransitionPage(
              child: const SigninScreen(key: Key('SigninScreen')),
              transitionDuration: const Duration(milliseconds: 150),
              transitionsBuilder: (context, animation, secondaryAnimation, child) => FadeTransition(opacity: animation, child: child),
            );
          },
        ),
        GoRoute(
          name: signupRouteName,
          path: signupRouteName,
          pageBuilder: (BuildContext context, GoRouterState state) {
            return CustomTransitionPage(
              child: const SignupScreen(key: Key('SignupScreen')),
              transitionDuration: const Duration(milliseconds: 150),
              transitionsBuilder: (context, animation, secondaryAnimation, child) => FadeTransition(opacity: animation, child: child),
            );
          },
        ),
        GoRoute(
          name: passwordRenewRouteName,
          path: passwordRenewRouteName,
          pageBuilder: (BuildContext context, GoRouterState state) {
            return CustomTransitionPage(
              child: const PasswordRenewScreen(key: Key('PasswordRenewScreen')),
              transitionDuration: const Duration(milliseconds: 150),
              transitionsBuilder: (context, animation, secondaryAnimation, child) => FadeTransition(opacity: animation, child: child),
            );
          },
        ),
        GoRoute(
          name: passwordResetRouteName,
          path: passwordResetRouteName,
          pageBuilder: (BuildContext context, GoRouterState state) {
            return CustomTransitionPage(
              child: const PasswordResetScreen(key: Key('PasswordResetScreen')),
              transitionDuration: const Duration(milliseconds: 150),
              transitionsBuilder: (context, animation, secondaryAnimation, child) => FadeTransition(opacity: animation, child: child),
            );
          },
        ),
        GoRoute(
          name: gameRouteName,
          path: gameRouteName,
          pageBuilder: (BuildContext context, GoRouterState state) {
            return CustomTransitionPage(
              child: const GameScreen(key: Key('GameScreen')),
              transitionDuration: const Duration(milliseconds: 150),
              transitionsBuilder: (context, animation, secondaryAnimation, child) => FadeTransition(opacity: animation, child: child),
            );
          },
        ),
        GoRoute(
          name: historyRouteName,
          path: historyRouteName,
          pageBuilder: (BuildContext context, GoRouterState state) {
            return CustomTransitionPage(
              child: const HistoryScreen(key: Key('HistoryScreen')),
              transitionDuration: const Duration(milliseconds: 150),
              transitionsBuilder: (context, animation, secondaryAnimation, child) => FadeTransition(opacity: animation, child: child),
            );
          },
        ),
        GoRoute(
          name: profilRouteName,
          path: profilRouteName,
          pageBuilder: (BuildContext context, GoRouterState state) {
            return CustomTransitionPage(
              child: const ProfilScreen(key: Key('ProfilScreen')),
              transitionDuration: const Duration(milliseconds: 150),
              transitionsBuilder: (context, animation, secondaryAnimation, child) => FadeTransition(opacity: animation, child: child),
            );
          },
        ),
        GoRoute(
          name: adminSigninRouteName,
          path: adminSigninRouteName,
          pageBuilder: (BuildContext context, GoRouterState state) {
            return CustomTransitionPage(
              child: const AdminSigninScreen(key: Key('AdminSigninScreen')),
              transitionDuration: const Duration(milliseconds: 150),
              transitionsBuilder: (context, animation, secondaryAnimation, child) => FadeTransition(opacity: animation, child: child),
            );
          },
        ),
        GoRoute(
          name: clientRouteName,
          path: clientRouteName,
          pageBuilder: (BuildContext context, GoRouterState state) {
            return CustomTransitionPage(
              child: const ClientScreen(key: Key('ClientScreen')),
              transitionDuration: const Duration(milliseconds: 150),
              transitionsBuilder: (context, animation, secondaryAnimation, child) => FadeTransition(opacity: animation, child: child),
            );
          },
        ),
      ],
    );
  }
}
