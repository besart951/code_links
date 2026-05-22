import type { AdminRole } from '$lib/domain/admin-access/types';
import type { UserListQuery, UserSortField, UserStatus } from '$lib/domain/users/types';

export class UserTableState {
	query = $state('');
	role = $state<AdminRole | 'user' | 'all'>('all');
	status = $state<UserStatus | 'all'>('all');
	page = $state(1);
	pageSize = $state(25);
	sort = $state<{ field: UserSortField; direction: 'asc' | 'desc' }>({
		field: 'createdAt',
		direction: 'desc'
	});

	request = $derived<UserListQuery>({
		query: this.query.trim() || undefined,
		role: this.role === 'all' ? undefined : this.role,
		status: this.status === 'all' ? undefined : this.status,
		page: this.page,
		pageSize: this.pageSize,
		sort: this.sort
	});

	setQuery = (query: string) => {
		this.query = query;
		this.page = 1;
	};

	setRole = (role: AdminRole | 'user' | 'all') => {
		this.role = role;
		this.page = 1;
	};

	setStatus = (status: UserStatus | 'all') => {
		this.status = status;
		this.page = 1;
	};
}
