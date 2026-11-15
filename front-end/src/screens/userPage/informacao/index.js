import React, { useState, useEffect } from "react";
import { View, Text, TouchableOpacity, TextInput, Image, ActivityIndicator, FlatList, Modal } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { useNavigation } from "@react-navigation/native";
import { Feather } from '@expo/vector-icons';

import styles from "./styled";
import { listarAnimais } from "../../../../services/API"; 

export default () => {
  const [animais, setAnimais] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [buscaTexto, setBuscaTexto] = useState("");
  
  const navigation = useNavigation();

  const carregarAnimais = async () => {
    try {
      setLoading(true);
      setError(null);
      const dados = await listarAnimais();
      setAnimais(dados);
    } catch (err) {
      console.error("Erro ao carregar os animais:", err);
      setError("Não foi possível carregar os animais.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    carregarAnimais();
  }, []);

  const animaisFiltrados = animais.filter((item) => {
    const termo = buscaTexto.toLowerCase();
    const nome = (item.nome || item.nome_popular || "").toLowerCase();
    const descricao = (item.descricao || item.especie || "").toLowerCase();

    return nome.includes(termo) || descricao.includes(termo);
  });

  const renderCardItem = ({ item }) => (
    <TouchableOpacity
      style={styles.cardContainer}
      onPress={() => navigation.navigate("animalDetalhe", { id: item.id, resumo: item })}
      activeOpacity={0.8}
    >
      <Image
        source={{ uri: item.imagem_url || "https://via.placeholder.com/300x180?text=Sem+Imagem" }}
        style={styles.cardImage}
        resizeMode="cover"
      />
      <View style={styles.cardTextArea}>
        <Text style={styles.cardTitle}>{item.nome || item.nome_popular || "Sem nome"}</Text>
        <Text style={styles.cardDescription} numberOfLines={3}>
          {item.breve_descricao || item.descricao || item.especie || "Sem descrição disponível."}
        </Text>
      </View>
    </TouchableOpacity>
  );

  return (
    <SafeAreaView style={styles.safeArea}>
      
      <View style={styles.header}>
        <View style={styles.headerSpacer} />
        <View style={styles.headerTitleBox}>
          <Text style={styles.headerTitleText}>ENCICLOPÉDIA</Text>
        </View>
        <View style={styles.headerSpacer} />
      </View>

      <View style={styles.buscaContainer}>
        <View style={styles.buscaInputWrapper}>
          <Feather name="search" size={20} color="#0A5133" style={styles.buscaIcon} />
          <TextInput 
            style={styles.buscaInput} 
            placeholder="Busca por nome ou espécie..." 
            placeholderTextColor="#5E6737"
            value={buscaTexto}                
            onChangeText={setBuscaTexto}      
            autoCapitalize="none"
          />
        </View>
      </View>

      {loading ? (
        <View style={styles.centerContainer}>
          <ActivityIndicator size="large" color="#A7F3D0" />
          <Text style={styles.loadingText}>Carregando animais...</Text>
        </View>
      ) : error ? (
        <View style={styles.centerContainer}>
          <Text style={styles.errorText}>{error}</Text>
          <TouchableOpacity style={styles.retryButton} onPress={carregarAnimais}>
            <Text style={styles.retryButtonText}>Tentar Novamente</Text>
          </TouchableOpacity>
        </View>
      ) : (
        <View style={{ flex: 1, width: "100%", marginTop: 10 }}>
          <FlatList
            data={animaisFiltrados}
            keyExtractor={(item, index) => item.id ? item.id.toString() : index.toString()}
            renderItem={renderCardItem}
            showsVerticalScrollIndicator={false}
            contentContainerStyle={{ paddingBottom: 20 }}
            ListEmptyComponent={
              <View style={styles.centerContainer}>
                <Text style={styles.emptyText}>Nenhum animal encontrado.</Text>
              </View>
            }
          />
        </View>
      )}
      
      <Modal 
        transparent={true}
        animationType='none'
        visible={loading}
      />
    </SafeAreaView>
  );
};