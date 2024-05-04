enum AuthStatus {
  connected("connected"),
  disconnected("disconnected");

  const AuthStatus(this.status);
  final String status;
}
