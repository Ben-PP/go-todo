class User {
  User({
    required this.id,
    required this.username,
    required this.isAdmin,
    required this.createdAt,
  });

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'] as String,
      username: json['username'] as String,
      isAdmin: json['is_admin'] as bool,
      createdAt: DateTime.parse(json['created_at'] as String),
    );
  }

  final String id;
  final String username;
  final bool isAdmin;
  final DateTime createdAt;

  @override
  String toString() {
    return 'User{id: $id, username: $username, isAdmin: $isAdmin, createdAt: $createdAt}';
  }
}
