class Todo {
  final String id;
  final String userId;
  final String listId;
  final String? parentId;
  final String title;
  final String description;
  final bool isCompleted;
  final DateTime createdAt;
  final DateTime updatedAt;
  final DateTime? completeBefore;
  final DateTime? completedAt;

  Todo({
    required this.id,
    required this.userId,
    required this.listId,
    this.parentId,
    required this.title,
    required this.description,
    required this.isCompleted,
    required this.createdAt,
    required this.updatedAt,
    this.completeBefore,
    this.completedAt,
  });

  factory Todo.fromJson(Map<String, dynamic> json) {
    return Todo(
      id: json['id'] as String,
      userId: json['user_id'] as String,
      listId: json['list_id'] as String,
      parentId: json['parent_id'] as String?,
      title: json['title'] as String,
      description: json['description'] as String,
      isCompleted: json['completed'] as bool,
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
      completeBefore: json['complete_before'] != null
          ? DateTime.parse(json['complete_before'] as String)
          : null,
      completedAt: json['completed_at'] != null
          ? DateTime.parse(json['completed_at'] as String)
          : null,
    );
  }

  @override
  String toString() {
    return 'Todo{id: $id, userId: $userId, listId: $listId, parentId: $parentId, title: $title, description: $description, isCompleted: $isCompleted, createdAt: $createdAt, updatedAt: $updatedAt, completeBefore: $completeBefore, completedAt: $completedAt}';
  }
}
