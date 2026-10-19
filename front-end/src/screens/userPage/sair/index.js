import React, { useEffect } from 'react';
import { View, ActivityIndicator, StyleSheet } from 'react-native';
import { useNavigation } from '@react-navigation/native';
import { sair } from '../../../../services/API';

export default () => {
    const navigation = useNavigation();

    useEffect(() => {
        const handleLogout = async () => {
            try {
                // 1. Chamar sua API de logout (se houver)git
                // await api.post('/logout');

                // 2. Limpar dados locais se necessário (ex: token)
                // await AsyncStorage.removeItem('@token');

                // 3. Voltar para o Stack Principal na tela de 'login', limpando o histórico
                navigation.reset({
                    index: 0,
                    routes: [{ name: 'login' }], // Nome da rota do Stack Navigator onde fica o login
                });
            } catch (error) {
                console.error("Erro ao fazer logout:", error);
                // Caso dê erro na API, ainda é seguro chutar o usuário para o login
                navigation.reset({
                    index: 0,
                    routes: [{ name: 'login' }],
                });
            }
        };

        handleLogout();
    }, [navigation]);

    // Um loading simples enquanto executa a limpeza/API
    return (
        <View style={styles.container}>
            <ActivityIndicator size="large" color="#2E7D32" />
        </View>
    );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
        backgroundColor: '#fff',
    },
});