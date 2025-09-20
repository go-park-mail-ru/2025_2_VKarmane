```mermaid
erDiagram 
    USER ||--o{ USER_ACCOUNT : "has" 
    USER { int id 
        string name 
        string surname 
        string email 
        string logo
        string login 
        string password 
        timestamptz created_at 
        timestamptz updated_at 
        } 


     ACCOUNT {
        int id
        decimal balance
        string type
        timestamptz created_at
        timestamptz updated_at
    }

    ACCOUNT ||--o{ USER_ACCOUNT : "belongs to" 
 
    USER_ACCOUNT{ 
        int id 
        int user_id 
        int account_id 
        timestamptz created_at
        timestamptz updated_at
        }

 
    OPERATION{
        int id
        int account_id
        int category_id
        string type
        string name
        decimal sum
        timestamptz created_at
    }

    CATEGORY ||--|{ OPERATION : "belongs to"
    USER ||--|{ CATEGORY : "has"
    CATEGORY {
        int id
        int user_id 
        string name
        string logo
        timestamptz created_at
        timestamptz updated_at
    }

    CURRENCY {
        int id
        string code
        string name
        string logo
        timestamptz created_at
    }

    CURRENCY ||--|{ACCOUNT : "currency of"
    CURRENCY ||--|{OPERATION : "currency of"

        TRANSFER {
        int id
        int from_account_id
    }
    ACCOUNT ||--o{ TRANSFER : "from"
    OPERATION ||--|{ ACCOUNT : "from/to"
    TRANSFER ||--|| OPERATION : "is"

BUDGET {
    int id
    int user_id
    decimal amount
    int currency_id
    string type
    bool is_failed
    timestamptz created_at
    timestamptz updated_at
    timestamptz closed_at
    date period_start
    date period_end
}


CURRENCY ||--o{ BUDGET : "currency of"
USER ||--|{ BUDGET : "has"