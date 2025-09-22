```mermaid
erDiagram 
    USER ||--o{ USER_ACCOUNT : "has" 
    USER { int id 
        string name 
        string surname 
        string email 
        int logo_id
        string login 
        string hashed_password 
        string user_description
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
        string status
        string description
        string receipt_url
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
        int logo_id
        string description
        timestamptz created_at
        timestamptz updated_at
    }

    CURRENCY {
        int id
        string code
        string name
        int logo_id
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
    string description
    timestamptz created_at
    timestamptz updated_at
    timestamptz closed_at
    date period_start
    date period_end
}

CHAT {
    int id
    string name
    timestamptz created_at
    timestamptz updated_at
}

USER_CHAT {
    int id
    user_id int
    chat_id int
    timestamptz created_at
    timestamptz updated_at
}

MESSAGE {
    int id
    int user_id
    int chat_id
    string message_text
    timestamptz created_at
    timestamptz updated_at
}

RECEIVER {
    int id
    int user_id
    string name
    timestamptz created_at
    timestamptz updated_at
}



CURRENCY ||--o{ BUDGET : "currency of"
USER ||--|{ BUDGET : "has"
USER ||-o{ MESSAGE : "writes"
USER ||--o{ USER_CHAT : "has"
USER ||--o{ RECEIVER : "has"
CHAT || -- { USER_CHAT : "connected to"