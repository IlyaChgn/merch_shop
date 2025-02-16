package repository

const (
	GetMerchItemQuery = `
		SELECT id, name, price
		FROM public.merchandise
		WHERE name = $1;
	`

	SaveItemQuery = `
		INSERT INTO public.inventory_item (user_id, merch_id) 
		VALUES ($1, $2); 
	`

	IncreaseBalanceQuery = `
		UPDATE public.balance
		SET amount = amount + $2
		WHERE user_id = $1;
	`

	DecreaseBalanceQuery = `
		UPDATE public.balance
		SET amount = amount - $2
		WHERE user_id = $1;
	`

	AddTransactionQuery = `
		INSERT INTO public.transaction (sender_id, receiver_id, amount)
		VALUES ($1, $2, $3);
	`

	GetBalanceQuery = `
		SELECT amount
		FROM public.balance
		WHERE user_id = $1;
	`

	GetIncomingTransactionsQuery = `
		SELECT t.amount, u.username
		FROM public.transaction t
		JOIN public.user u
			ON t.sender_id = u.id
		WHERE t.receiver_id = $1;
	`

	GetOutgoingTransactionsQuery = `
		SELECT t.amount, u.username
		FROM public.transaction t
		JOIN public.user u
			ON t.receiver_id = u.id
		WHERE t.sender_id = $1;
	`

	GetInventoryQuery = `
		SELECT m.name, COUNT(i.*)
		FROM public.inventory_item i
		JOIN public.merchandise m
			ON i.merch_id = m.id
		WHERE i.user_id = $1
		GROUP BY m.name;
	`
)
