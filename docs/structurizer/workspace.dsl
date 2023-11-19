workspace "Анти-брутфорс" {
    !docs docs

    model {
        AntiBruteforce = softwareSystem "Анти-брутфорс" "Сервис предназначен для борьбы с подбором паролей при авторизации в какой-либо системе."
        InMemoryCache = softwareSystem "In memory DB" "" "Database"
        PostgreSQL = softwareSystem "PostgreSQL Database" "" "Database"

        Client = softwareSystem "Client" "" "External System"

        RateLimiter = softwareSystem "Rate Limiter" "" "Internal System" {
            API = container "API Rate Limiter" "Go" "Container" {
                AuthController = component "AuthController" "" "Go" "Component"
                BlacklistController = component "BlacklistController" "" "Go" "Component"
                WhiteListController = component "WhiteListController" "" "Go" "Component"
                BucketController = component "BucketController" "" "Go" "Component"

                AuthUseCase = component "AuthUseCase" "" "Go" "Component,UseCase"
                BlacklisthUseCase = component "BlacklisthUseCase" "" "Go" "Component,UseCase"
                WhiteListUseCase = component "WhiteListUseCase" "" "Go" "Component,UseCase"
                BucketUseCase = component "BucketUseCase" "" "Go" "Component,UseCase"

                RateLimiterService = component "RateLimiterService" "" "Go" "Compoenent,Service"
                BlacklisthService = component "BlacklisthService" "" "Go" "Component,Service"
                WhiteListService = component "WhiteListService" "" "Go" "Component,Service"
                BucketService = component "BucketService" "" "Go" "Component,Service"

                ListRepository = component "ListRepository" "" "Go" "Compoenent,Repository"

                BucketRepository = component "BucketRepository" "" "Go" "Compoenent,Repository"
                IpBucketRepository = component "IpBucketRepository" "" "Go" "Compoenent,Repository"
                LoginBucketRepo = component "LoginBucketRepo" "" "Go" "Compoenent,Repository"
                PasswordBucketRepo = component "PasswordBucketRepo" "" "Go" "Compoenent,Repository"
            }
        }

        Client -> RateLimiter "" "JSON/HTTP"

        AuthController -> AuthUseCase "Use"
        BlacklistController -> BlacklisthUseCase "Use"
        WhiteListController -> WhiteListUseCase "Use"
        BucketController -> BucketUseCase "Use"

        AuthUseCase -> BlacklisthService "Use"
        AuthUseCase -> WhiteListService "Use"
        AuthUseCase -> BucketService "Use"

        BlacklisthUseCase -> BlacklisthService "Use"
        WhiteListUseCase -> WhiteListService "Use"
        BucketUseCase -> BucketService "Use"

        BlacklisthService -> ListRepository "Use"
        WhiteListService -> ListRepository "Use"

        BucketService -> IpBucketRepository "Use"
        BucketService -> LoginBucketRepo "Use"
        BucketService -> PasswordBucketRepo "Use"

        IpBucketRepository -> BucketRepository "Use"
        LoginBucketRepo -> BucketRepository "Use"
        PasswordBucketRepo -> BucketRepository "Use"

        BucketRepository -> InMemoryCache "Use"
        ListRepository -> PostgreSQL "Use"
    }

    views {
        component API "Overview" {
            include *
            autoLayout
        }

        theme default
        styles {
            element "Element" {
                shape RoundedBox
            }
            element "Person" {
                shape Person
            }
            element "External System" {
                background #999999
                color #000000
            }
            element "Redis" {
                background #1abc9c
                color #ffffff
            }
            element "Internal System" {
                background #6600cc
                color #ffffff
            }
            element "Software System" {
                background #1168bd
                color #ffffff
            }
            element "Container" {
                background #438dd5
                color #ffffff
            }
            element "Component" {
                background #85bbf0
                color #000000
            }
            element "Database" {
                shape Cylinder
            }
            element "Kafka" {
                shape Pipe
            }
        }
    }
}

